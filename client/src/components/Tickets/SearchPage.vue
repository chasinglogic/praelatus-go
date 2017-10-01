<template>
  <div class="container">
    <h1>Tickets</h1>
    <b-form-fieldset class="mr-auto ml-auto">
      <b-form-input v-model="query" @keyup="loadTickets"
        placeholder="Type to Search" />
    </b-form-fieldset>
    <ticket-list :tickets="tickets" showColumnPicker="true"></ticket-list>
  </div>
</template>

<script>
 import TicketList from '@/components/Tickets/List'
 import Axios from 'axios'

 export default {
   components: {
     TicketList
   },

   data: () => {
     return {
       'query': '',
       'tickets': []
     }
   },

   methods: {
     loadTickets: function () {
       let url = '/api/tickets'
       let inst = this

       if (this.query && this.query !== '') {
         url += '?q=' + this.query
       }

       Axios.get(url)
            .then((res) => {
              console.log(res.data)
              inst.tickets = res.data
            })
            .catch((err) => {
              if (err.response.status === 404) {
                inst.tickets = []
                return
              }

              // TODO: Visually que the user that there's been an error
              console.log('ERROR', err)
            })
     }
   },

   created: function () {
     this.loadTickets()
   }
 }
</script>
