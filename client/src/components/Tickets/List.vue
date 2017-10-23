<!-- Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights
     reserved. Use of this source code is governed by the AGPLv3 license that
     can be found in the LICENSE file. -->

<template>
  <div id="ticket-list-root">
    <div id="list-wrapper" v-if="tickets.length !== 0">
      <table class="table">
        <thead>
          <tr>
            <template v-for="column in columns">
              <th v-show="column.active">{{ column.displayName ? column.displayName : humanizeColumnName(column.name) }}</th>
            </template>
            <th v-if="showColumnPicker">
              <b-dropdown text="Columns">
                <div v-for="column in columns">
                  <span>{{ column.displayName ? column.displayName : humanizeColumnName(column.name) }}</span>
                  <input type="checkbox" v-model="column.active" />
                </div>
                <div>
                  <b-button @click="resetDefaultColumns">
                    Reset Defaults
                  </b-button>
                </div>
              </b-dropdown>
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="ticket in tickets">
            <template v-for="column in columns">
              <td v-show="column.active" v-if="column.name === 'key'">
                <router-link v-bind:to="'/tickets/' + ticket.key">
                  {{ ticket.key }}
                </router-link>
              </td>
              <td v-show="column.active" v-else-if="ticket[column.name]">
                {{ ticket[column.name] }}
              </td>
              <td v-show="column.active" v-else>
                {{ getFieldValue(ticket, column.name) }}
              </td>
            </template>
          </tr>
        </tbody>
      </table>
    </div>
    <div v-else>
      <loading-spinner></loading-spinner>
    </div>
  </div>
</template>

<script>
 import LoadingSpinner from '@/components/General/LoadingSpinner'

 export default {
   name: 'ticket-list',
   components: {
     LoadingSpinner
   },
   methods: {
     resetDefaultColumns: function () {
       this.columns = this.defaultColumns()
     },

     getFieldValue: function (ticket, fieldName) {
       let field = ticket.fields.filter(f => f.name === fieldName)
       return field ? field.value : 'None'
     },

     humanizeColumnName: function (columnName) {
       return columnName
       // insert a space before all caps
         .replace(/([A-Z])/g, ' $1')
       // uppercase the first character
         .replace(/^./, function (str) { return str.toUpperCase() })
         .replace(/^ /, '')
         .replace('\n', '')
     }
   },

   watch: {
     tickets: function () {
       if (this.tickets[0]) {
         return this.columns
                    .concat(
                      this
                        .tickets[0]
                        .fields
                        .map(f => { return { name: f.name, active: true } })
                    )
       }

       this.columns = Array.from(this.defaultColumns)
     }
   },

   data: function () {
     let defaults = () => {
       return [
         {
           name: 'key',
           active: true
         },
         {
           name: 'summary',
           active: true
         },
         {
           name: 'createdDate',
           active: false
         },
         {
           name: 'updatedDate',
           active: false
         },
         {
           name: 'status',
           active: true
         },
         {
           name: 'project',
           active: true
         },
         {
           name: 'description',
           active: false
         },
         {
           name: 'assignee',
           active: false
         },
         {
           name: 'reporter',
           active: false
         },
         {
           name: 'labels',
           active: false
         },
         {
           name: 'ticketType',
           active: true
         }
       ]
     }
     return {
       defaultColumns: defaults,
       columns: defaults()
     }
   },

   props: {
     'showColumnPicker': false,
     'tickets': []
   }
 }
</script>


<style>
 th {
   text-align: center;
 }
</style>
