<!-- Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights
     reserved. Use of this source code is governed by the AGPLv3 license that
     can be found in the LICENSE file. -->

<template>
  <div class="card">
    <h2 class="card-header" >
      Details
    </h2>
    <div class="card-block" >
      <div class="container">
        <p>
          <span class="field-name">Status:</span>
        </p>
        <status-pill :status="ticket.status" />
        <br />

        <p>
          <span class="field-name">Reporter:</span>
          <user-stub :username="ticket.reporter" />
        </p>

        <p>
          <span class="field-name">Assignee:</span>
          <user-stub :username="ticket.assignee" />
        </p>

        <p>
          <span class="field-name">Created:</span>
          {{ dateFormat(ticket.createdDate) }}
        </p>

        <p>
          <span class="field-name">Updated:</span>
          {{ dateFormat(ticket.updatedDate) }}
        </p>

        <p>
          <span class="field-name">Type:</span>
          {{ ticket.type }}
        </p>

        <template v-if="ticket.labels.length !== 0">
          <p>
            <span class="field-name">Labels:</span>
          </p>
          <div class="label" v-for="label in ticket.labels">
            {{ label }}
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<script>
 import UserStub from '@/components/Users/Stub'
 import StatusPill from '@/components/Status/Pill'
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
       default: {
         'labels': []
       }
     }
   }
 }
</script>


<style>
 .label {
   background-color: #eee;
   color: blue;
   border-radius: 5px;
   margin-left: 0.5rem;
   margin-bottom: 0.5rem;
   padding: 0.2rem;
   text-align: center;
   display: inline-block;
 }

 .field-name {
   font-weight: bold;
   margin-right: 0.3rem;
 }
</style>
