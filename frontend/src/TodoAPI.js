import axios from 'axios';
import  BACKEND_URL from "./config";
const baseUrl = BACKEND_URL

console.log(BACKEND_URL )
export default class TodoApi{
     static async getAll(){
          try {
               const res = await axios.get(baseUrl+"/api/todo/all")
               return res.data
          }catch(err) {
               console.log(err)
          }
     }
     static async updateTodo(id,task,done) {
          await axios.put(baseUrl +"/api/todo/"+"?taskId="+id,{
               task:task,
               done:done
          }).then(function(res) {
               console.log(res)
          }).catch((err)=>{
               console.log(err)
          })
     }
     static async deleteTodo(id) {
          await axios.delete(baseUrl +"/api/todo/"+"?taskId="+id).then(function(res) {
               console.log(res)
          }).catch((err)=>{
               console.log(err)
          })
     }
     static async newTodo(task) {
          await axios.post(baseUrl +"/api/todo",{
               task:task,
          }).then(function(res) {
               console.log(res)
          }).catch((err)=>{
               console.log(err)
          })
     }
}