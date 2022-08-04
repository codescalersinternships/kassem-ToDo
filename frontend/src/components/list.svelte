<script>
     import {onMount} from "svelte"
     import TodoApi from "../TodoAPI"
     import {items} from "../store"
     import Item from "./item.svelte"
     import NewItem from "./newItem.svelte"
     function handleNewItem(event) {
          TodoApi.newTodo(event.detail)
          .then( async() =>$items = await TodoApi.getAll())
     }
     function handleUpdateItem(event) {

         console.log("update with",event.detail.ID,event.detail.task, event.detail.done)
          TodoApi.updateTodo(event.detail.ID,event.detail.task, event.detail.done)
          .then( async() =>$items = await TodoApi.getAll())
     }
     let i =0
     onMount(async()=> {
          // fetch from AP
          $items = await TodoApi.getAll()
     })
</script>
   
<style>
     .list{
          padding:15px;
     }
     .list-status{
          margin:0;
          text-align:center;
          color: #ffffff;
          font-weight: bold;
          font-size: 1.1em;
     }
</style>

<div class="list">
     <NewItem on:newItem={handleNewItem}/>
    
     {#each $items as item ,i (item)}
  
      <Item {...item} clsss="hhh" on:update={handleUpdateItem}  />
         
     {:else}
          <p class='list-status'>No Items exist</p>
     {/each}
</div>