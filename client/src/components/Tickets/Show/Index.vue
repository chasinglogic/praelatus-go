<!-- Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights
     reserved. Use of this source code is governed by the AGPLv3 license that
     can be found in the LICENSE file. -->

<template>
  <div v-if="loading">
    <h1>Loading your ticket...</h1>
  </div>
  <div v-else class="container-fluid ticket-layout">
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
          <comment-form />
          <comment-form v-if="currentUser" />
          <div class="card comment" v-else>
            <div class="comment card-block">
              You must be <a href="/#/login">logged in</a> to comment.
            </div>
          </div>
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
 import BreadCrumb from './Children/BreadCrumb'
 import TicketDetails from './Children/Details'
 import TicketFields from './Children/Fields'

 import Comments from '@/components/Comments/List'
 import Sidebar from '@/components/General/Sidebar'

 import Axios from 'axios'
 import Showdown from 'showdown'
 const converter = new Showdown.Converter()

 export default {
   name: 'ticket',
   components: {
     Sidebar,
     BreadCrumb,
     Comments,
     'ticket-fields': TicketFields,
     'ticket-details': TicketDetails
   },

   data: () => {
     return {
       'loading': true,
       'ticket': {
         'labels': [],
         'fields': [],
         'comments': []
       }
     }
   },

   methods: {
     loadTicket: function () {
       let url = '/api/tickets/' + this.$route.params.key
       let inst = this

       Axios.get(url)
            .then((res) => {
              inst.ticket = res.data
              inst.loading = false
            })
            .catch((err) => {
              if (err.response.status === 404) {
                this.$route.router.go('*')
              }

              console.log('ERROR', err.response.data)
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
