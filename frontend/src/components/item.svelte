<script>
     import {createEventDispatcher} from "svelte"
     export let ID, task, done;
     import {items} from "../store"
     import TodoApi from "../TodoAPI"
     const dispatch = createEventDispatcher()
     function triggerUpdate() {
          console.log(done)
          dispatch("update",{ID, task, done})
     }
     function handlebtnClick(event){
          TodoApi.deleteTodo(event.target.id)
          .then( async() =>$items = await TodoApi.getAll())
     }
     
     
</script>
<style>
.item {
     display: flex;
     align-items: center;
     padding: 15px;
     background : #ffffff;
     
}
.item:hover{
     background : rgba(255, 255, 255,0.8)
}
.item.done{
     background :#dddddd;
     
}
.item.done .text-input{
     color: #555555;
     text-decoration: line-through;
}
.text-input {
     flex-grow:1;
     background: none;
     border: none;
     outline:none;
     font-weight: 500;
}
.completed-checkbox {
     margin-left: 15px;
     margin-top:1px
}
.delete{
     font-size: 12px;
     margin:1px;
     padding: 5px;
}

</style>
<div
 class="{done == true ? " item done" : "item"}" 
 >
 
     <input class="completed-checkbox"
      type="checkbox"
      bind:checked={done}
      on:change={()=> triggerUpdate()}
      />
     <input 
          type="text"
          class="text-input" 
          bind:value={task} 
          readonly={done}
          on:keyup={({key, target}) => key==="Enter" && target.blur()}
          on:blur = {() => triggerUpdate()}
     />
     <button class="btn btn-danger delete" id={ID}
     
      on:click={handlebtnClick}
      
      >x
     </button>
     
</div>