<template>
  <div class="container-fluid ticket-layout">
    <sidebar></sidebar>
    <div class="container-fluid">
      <div class="ticket-header card">
        <h3 class="card-header">
            <bread-crumb :ticket="ticket" />
        </h3>
        <h1 class="card-block">
          {{ ticket.summary }}
        </h1>
      </div>
      <div class="row" >
        <div class="col-md-9">
          <div class="card">
            <h2 class="card-header">
              Description
            </h2>
            <div v-html="markdown(ticket.description)" class="card-block">
            </div>
          </div>
          <comments :comments="ticket.comments" />
        </div>
        <div class="col-md-3" >
          <ticket-details :ticket="ticket" />
          <ticket-fields :fields="ticket.fields" />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
 import BreadCrumb from './BreadCrumb'
 import Comments from './Comments'
 import TicketDetails from './Details'
 import Sidebar from '@/components/General/Sidebar'

 import Showdown from 'showdown'
 const converter = new Showdown.Converter()

 export default {
   name: 'ticket',
   components: {
     Sidebar,
     BreadCrumb,
     'ticket-details': TicketDetails,
     Comments
   },

   computed: {
     ticket: function () {
       return this.$store.getters.currentTicket
     }
   },

   methods: {
     loadTicket: function () {
       let url = '/api/tickets/' + this.$route.params.key
       this.$store.dispatch('request', {
         url: url,
         key: 'currentTicket'
       })
     },

     markdown (text) {
       return converter.makeHtml(text)
     }
   },

   created: function () {
     this.loadTicket()
   }
 }
</script>

<style>
 .ticket-layout {
   text-align: left;
   padding-top: 1rem;
 }

 .ticket-header {
   margin-bottom: 1rem;
 }
</style>
