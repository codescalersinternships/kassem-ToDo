import axios from 'axios';

const baseUrl = BACKEND_URL
console.log(baseUrl )
export default class TodoApi{
     static async getAll(){
          try {
               const res = await axios.get(baseUrl+"todo/all")
               return res.data
          }catch(err) {
               console.log(err)
          }
     }
     static async updateTodo(id,task,done) {
          await axios.put(baseUrl +"todo/"+"?taskId="+id,{
               task:task,
               done:done
          }).then(function(res) {
               console.log(res)
          }).catch((err)=>{
               console.log(err)
          })
     }
     static async deleteTodo(id) {
          await axios.delete(baseUrl +"todo/"+"?taskId="+id).then(function(res) {
               console.log(res)
          }).catch((err)=>{
               console.log(err)
          })
     }
     static async newTodo(task) {
          await axios.post(baseUrl +"todo",{
               task:task,
          }).then(function(res) {
               console.log(res)
          }).catch((err)=>{
               console.log(err)
          })
     }
}