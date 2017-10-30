<!-- Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights
     reserved. Use of this source code is governed by the AGPLv3 license that
     can be found in the LICENSE file. -->

<template>
  <div id="item-list-root">
    <div id="list-wrapper" v-if="items.length !== 0">
      <table class="table">
        <thead>
          <tr>
            <template v-for="column in columns">
              <th>{{ humanizeColumnName(column) }}</th>
            </template>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in items">
            <template v-for="column in columns">
              <td v-if="column === uid">
                <router-link v-bind:to="prefix + item[column]">
                  {{ item[column] }}
                </router-link>
              </td>
              <td v-else-if="item[column]">
                {{ item[column] }}
              </td>
              <td v-else>
                None
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
   name: 'item-list',
   components: {
     LoadingSpinner
   },

   methods: {
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
     items: function () {
       if (this.items[0] && this.columns.length === 0) {
         this.columns = Object.keys(this.items[0])
         this.columns = this.columns.filter(this.ignoredColumns)
       }
     }
   },

   props: {
     'columns': {
       type: Array,
       default: []
     },
     'items': {
       type: Array,
       default: []
     },
     'uid': {
       type: String,
       default: ''
     },
     'ignoredColumns': {
       type: Function,
       default: function (col) { return true }
     },
     'prefix': {
       type: String,
       default: ''
     }
   }
 }
</script>


<style>
 th {
   text-align: center;
 }
</style>
