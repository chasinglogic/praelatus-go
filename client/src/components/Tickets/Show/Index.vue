<!-- Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights
     reserved. Use of this source code is governed by the AGPLv3 license that
     can be found in the LICENSE file. -->

<template>
  <div>
    <div v-if="loading" class="container">
      <loading-spinner></loading-spinner>
    </div>
    <div v-else class="container ticket-layout">
      <div class="ticket-header card">
        <h3 class="card-header">
          <bread-crumb :ticket="ticket" />
        </h3>
        <h1 class="card-block">
          {{ ticket.summary }}
        </h1>
      </div>
      <div class="row" >
        <div class="col-lg-8 col-12">
          <div class="card">
            <h2 class="card-header">
              Description
            </h2>
            <div v-html="markdown(ticket.description)" class="card-block">
            </div>
          </div>
        </div>
        <div class="col-lg-4 col-12" >
          <div class="row">
            <div class="col-6-lg col-12">
              <ticket-details :ticket="ticket" />
            </div>
            <div class="col-6-lg col-12">
              <ticket-fields :fields="ticket.fields" />
            </div>
          </div>
        </div>
      </div>
      <comments :comments="ticket.comments" />
      <comment-form @newComment="loadTicket" />
    </div>
  </div>
</template>

<script>
 import BreadCrumb from './Children/BreadCrumb'
 import TicketDetails from './Children/Details'
 import TicketFields from './Children/Fields'

 import Comments from '@/components/Comments/List'
 import CommentForm from '@/components/Comments/Form'
 import Sidebar from '@/components/General/Sidebar'
 import LoadingSpinner from '@/components/General/LoadingSpinner'

 import Markdown from '@/lib/markdown'
 import Axios from 'axios'

 export default {
   name: 'ticket',
   components: {
     Sidebar,
     BreadCrumb,
     Comments,
     LoadingSpinner,
     'comment-form': CommentForm,
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
                this.$router.push('/404')
              }

              console.log('ERROR', err.response.data)
            })
     },

     markdown: Markdown.render
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

 .card-block {
   padding: 1rem;
 }
</style>
