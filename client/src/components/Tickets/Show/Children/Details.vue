<!-- Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights
     reserved. Use of this source code is governed by the AGPLv3 license that
     can be found in the LICENSE file. -->

<template>
  <div class="card">
    <h2 class="card-header" >
      Details
    </h2>
    <div class="card-block" >
      <div class="container">
        <div class="field-name">
          Status:
        </div>
        <div class="field-value">
          <status-pill :status="ticket.status" />
        </div>

        <div class="field-name">
          Reporter:
        </div>
        <div class="field-value">
          <user-stub :username="ticket.reporter" />
        </div>

        <div class="field-name">
          Assignee:
        </div>
        <div class="field-value">
          <template v-if="ticket.assignee">
            <user-stub :username="ticket.assignee" />
          </template>
          <template v-else>
            Unassigned
          </template>
        </div>

        <div class="field-name">
          Created:
        </div>
        <div class="field-value">
          {{ dateFormat(ticket.createdDate) }}
        </div>

        <div class="field-name">
          Updated:
        </div>
        <div class="field-value">
          {{ dateFormat(ticket.updatedDate) }}
        </div>

        <div class="field-name">
          Type:
        </div>
        <div class="field-value">
          {{ ticket.type }}
        </div>

        <template v-if="ticket.labels.length !== 0">
          <div class="field-name">
            Labels:
          </div>
          <div class="field-value">
            <div class="label" v-for="label in ticket.labels">
              <router-link :to="{ path: '/queries',
                                query: { q: labelQuery(label)} }">
                {{ label }}
              </router-link>
            </div>
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
     dateFormat: dateUtils.dateFormat,

     labelQuery: function (label) {
       return 'labels = "' + label + '"'
     }
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

 .field-value .status-pill {
   padding: 0.5rem;
 }
</style>
