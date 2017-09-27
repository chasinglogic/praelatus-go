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
 import SearchBar from '@/components/General/SearchBar'
 import Axios from 'axios'

 export default {
   components: {
     SearchBar,
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
              inst.tickets = res.data
            })
            .catch((err) => {
              console.log('ERROR', err)
            })
     }
   },

   created: function () {
     this.loadTickets()
   }
 }
</script>
