 
it('open website', () => {
 cy.visit('http://localhost:8080')
})
it('check for headers', () => {
  cy.get('.header > :nth-child(1)')
})
it("check for counter", () => {
  cy.get('.header > :nth-child(2)')
})

it("clean database",()=>{
  cy.get(".header > :nth-child(2)").then(function($elem) {
    let countArr=($elem.text()).split("/")
    
    let index= parseInt(countArr[1])
    console.log(index)
    if(index==0){
      cy.get('.list-status')
    }
  })
})

it("add new task", () => {
  cy.get("input[placeholder='+ Add a task']")
  .type('please{enter}')
})

it('remove task', () => {
  cy.get("body > div:nth-child(2) > div:nth-child(2) > div:nth-child(2) > button:nth-child(3)")
  .click()  
}) 

it("check for done",() => {
  cy.get("input[placeholder='+ Add a task']")
  .type('test to be checked{enter}')

  cy.get(":nth-child(2) > .completed-checkbox")
  .check()  
})

it("check for the counter",() => {
  cy.get(".header > :nth-child(2)").then(function($elem) {
    let countArr=($elem.text()).split("/")
    
    let index= parseInt(countArr[0])
    console.log(index)
    if(index>0){
      return true
    }
  })

})
it("remove task after completed",()=>{
  cy.get("body > div:nth-child(2) > div:nth-child(2) > div:nth-child(2) > button:nth-child(3)")
  .click()  
})

it("add task to be updated",()=>{
  cy.get("input[placeholder='+ Add a task']")
  .type('test to be updated{enter}')
})
it("update task",() => {

  cy.get("body > div:nth-child(2) > div:nth-child(2) > div:nth-child(2)").find('.text-input')
  .type('..... ok its updated{enter}')
})

it("remove task after updated",()=>{
  cy.get("body > div:nth-child(2) > div:nth-child(2) > div:nth-child(2) > button:nth-child(3)")
  .click()  
})
