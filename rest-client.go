package main

import "io/ioutil"
import "fmt"
import "bufio"
import "os"
import "strings"
import "net/http"
import "errors"

var stayAlive bool = true;
var myName = "";
var site = "";

//COMMANDS
const COMMAND_PREFIX string = "/";
const HELP_COMMAND string = COMMAND_PREFIX+"help";
const QUIT_COMMAND string = COMMAND_PREFIX+"quit";
const CREATE_ROOM_COMMAND string = COMMAND_PREFIX+"createRoom"; //creates a room with the name of the first argument given
const LIST_ROOMS_COMMAND string = COMMAND_PREFIX+"listRooms"
const JOIN_ROOM_COMMAND string = COMMAND_PREFIX+"join";//   /join roomname will add a user to a rooms list of clients and switch the user to that room
const CURR_ROOM_COMMAND string = COMMAND_PREFIX+"currentRoom";
const CURR_ROOM_USERS_COMMAND string = COMMAND_PREFIX+"currentUsers";
const LEAVE_ROOM_COMMAND string = COMMAND_PREFIX+"leaveRoom";

//continusly asks the server for input by calling for messages for the user
func getFromServer(){
  for stayAlive {
    resp, err := http.Get(site+"/"+myName+"/messages")
    if err != nil{
      fmt.Println("error in getting messages")
      fmt.Println(err)
      stayAlive = false;
      return
    }
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Print(string(body))
  }
}


type DoubleArgs struct{
  Arg1 string;
  Arg2 string;
}

//creates the http message that will be sent to the server
func messageHelper(method string, url string) error{
  client := &http.Client{
    CheckRedirect: nil,
  }
    reply, err  := http.NewRequest(method, url, nil)
    reply.Header.Add("username", myName)
    client.Do(reply)
    return err

}

//Handles user input, reads from stdin and then posts that line to the server
func getfromUser(){

    for stayAlive{
      reader := bufio.NewReader(os.Stdin)
      message, _ := reader.ReadString('\n')//read from stdin till the next newline
      var err error;
      message = strings.TrimSpace(message);//strips the newlines from the input
      isCommand := strings.HasPrefix(message, COMMAND_PREFIX);//checks to see if the line starts with /
      if(isCommand){
        //parse command line, commands should be in the exact form of "/command arg arg arg" where args are not required
        parsedCommand := strings.Split(message, " ")
        if parsedCommand[0] == HELP_COMMAND {
          err = messageHelper("GET", site+"/help")
        } else if parsedCommand[0] == QUIT_COMMAND {
          err = messageHelper("DELETE", site+"/"+myName)
          stayAlive = false;
        } else if parsedCommand[0] == CREATE_ROOM_COMMAND {
          // not enough arguments to the command
          if len(parsedCommand) < 2{
            err = errors.New("not enough args for create room")
          }else{
            err = messageHelper("POST", site+"/rooms/"+parsedCommand[1])
          }
        } else if parsedCommand[0] == LIST_ROOMS_COMMAND {
          err = messageHelper("GET", site+"/rooms")
        } else if parsedCommand[0] == JOIN_ROOM_COMMAND {
          //not enough given to the command
          if len(parsedCommand) < 2{
            err = errors.New("You must specify a room to join")
          }else{
            err = messageHelper("POST", site+"/rooms/"+parsedCommand[1]+"/"+myName)
          }
        } else if parsedCommand[0] == CURR_ROOM_COMMAND {
          err = messageHelper("GET", site+"/"+myName+"/currentroom")
        }else if parsedCommand[0] == CURR_ROOM_USERS_COMMAND{
          err = messageHelper("GET", site+"/"+myName+"/currentroomusers")
        }else if parsedCommand[0] == LEAVE_ROOM_COMMAND{
          err = messageHelper("DELETE", site+"/"+myName+"/leaveroom")
        }

      }else if stayAlive{ // message is not a command
        //we need to create a post request to send the message to the server
        client := &http.Client{
          CheckRedirect: nil,
        }
          sendReply, _  := http.NewRequest("POST", site+"/"+myName+"/messageRoom", nil)
          sendReply.Header.Add("message", message)
          client.Do(sendReply)

      }
      if err != nil{
        fmt.Println(err)
      }
    }
  }

//starts up the client, starts the recieving thread and the input threads and then loops forever
func main() {

arguments := os.Args[1:];
IP := "localhost";
PORT:= "8080";
if len(arguments) == 0 {
  //no arguments start on localhost 8080
} else if len(arguments) != 2 {
  fmt.Println("I cannot understand your arguments, you must specify no arguments or exactly 2, first the IP and the second as the port")
  return
} else if len(arguments) == 2 {
//correct ammount of args
IP = arguments[0]
PORT = arguments[1]
}
//fmt.Println(arg)
  // connect to this socket
  fmt.Println("Attempting to connect to "+IP+":"+PORT)
  site = "http://"+IP+":"+PORT;
  resp, err := http.Get(site)
  if err != nil{
    fmt.Println("Something went wrong with the connection, check that the server exists and that your IP/Port are correct:\nError Message: ")
    fmt.Println(err)
    return
  }
  defer resp.Body.Close()
  body, _ := ioutil.ReadAll(resp.Body)
  fmt.Println(string(body))
  myName = string(body)
  fmt.Println("Your Username is: "+myName)
  go getfromUser();
  go getFromServer()
  for stayAlive {
    //loops  forever until stayAlive is set to false and then it shuts down
  }
}
