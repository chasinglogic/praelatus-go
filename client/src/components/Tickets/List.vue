<!-- Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights
     reserved. Use of this source code is governed by the AGPLv3 license that
     can be found in the LICENSE file. -->

<template>
  <div id="ticket-list-root">
    <div id="list-wrapper" v-if="tickets !== null">
      <table class="table">
        <thead>
          <tr>
            <template v-for="column in columns">
              <th v-show="column.active">{{ column.displayName ? column.displayName : humanizeColumnName(column.name) }}</th>
            </template>
            <th v-if="showColumnPicker">
              <b-dropdown text="Columns" variant="outline-primary"
                class="column-dropdown">
                <div v-for="column in columns">
                  <div class="column-dropdown-name">
                    {{ humanizeColumnName(column.name) }}
                  </div>
                  <input type="checkbox" class="column-dropdown-input"
                    v-model="column.active" />
                </div>
                <b-button class="column-dropdown-reset-button"
                  variant="warning"
                  @click="resetDefaultColumns">
                  Reset Defaults
                </b-button>
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
       this.columns = this.defaults()
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
     tickets: function (newVal) {
       console.log(newVal)
       if (newVal && newVal.length !== 0) {
         return this.columns
                    .concat(
                      newVal[0]
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
       defaults: defaults,
       defaultColumns: defaults(),
       columns: defaults()
     }
   },

   props: {
     'showColumnPicker': false,
     'tickets': {
       type: Array,
       default: null
     }
   }
 }
</script>


<style>
 th {
   text-align: center;
 }

 .column-dropdown-name {
   display: inline-block;
   margin-left: 0.5rem;
 }

 .column-dropdown-input {
   float: right;
   margin-right: 0.5rem;
 }

 .column-dropdown-reset-button {
   width: 100%;
   margin-top: 0.5rem;
 }

 .column-dropdown .dropdown-menu {
   padding-bottom: 0;
 }
</style>
