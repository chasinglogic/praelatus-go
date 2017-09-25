<template>
  <div class="card details">
    <h2 class="card-header" >
      Details
    </h2>
    <div class="card-block" >
      <p>
        <span class="bold">Status:</span>
      </p>
      <status-pill :status="ticket.status" />
      <br />

      <p>
        <span class="bold">Reporter:</span>
        <user-stub :username="ticket.reporter" />
      </p>

      <p>
        <span class="bold">Assignee:</span>
        <user-stub :username="ticket.assignee" />
      </p>

      <p>
        <span class="bold">Created:</span>
        {{ dateFormat(ticket.createdDate) }}
      </p>

      <p>
        <span class="bold">Updated:</span>
        {{ dateFormat(ticket.updatedDate) }}
      </p>

      <p>
        <span class="bold">Type:</span>
        {{ ticket.ticketType }}
      </p>

      <template v-if="ticket.labels.length !== 0">
        <p>
          <span class="bold">Labels:</span>
        </p>
        <p>
          <span v-for="label in labels">{{ label }}</span>
        </p>
      </template>
    </div>
  </div>
</template>

<script>
 import UserStub from '../Users/UserStub'
 import StatusPill from './StatusPill'
 import dateUtils from '@/lib/dates'

 export default {
   name: 'ticket-details',

   components: {
     'user-stub': UserStub,
     'status-pill': StatusPill
   },

   methods: {
     dateFormat: dateUtils.dateFormat
   },

   props: {
     ticket: {
       name: 'ticket',
       default: {}
     }
   }
 }
</script>


<style>
 .detail p {
   display: inline-block;
 }

 .bold {
   font-weight: bold;
 }

 .avatar {
   width: 20px;
   height: 20px;
   display: inline-block
 }
</style>
