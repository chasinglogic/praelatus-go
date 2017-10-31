<!-- Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights
     reserved. Use of this source code is governed by the AGPLv3 license that
     can be found in the LICENSE file. -->

<template>
  <div class="card fields">
    <h2 class="card-header" >
      Fields
    </h2>

    <div class="card-block" >
      <div v-for="field in fields">
        <div class="field-name">{{ field.name }}:</div>
        <div class="field-value">{{ getValue(field) }}</div>
      </div>
    </div>
  </div>
</template>

<script>
 import dateUtils from '@/lib/dates'

 export default {
   name: 'ticket-fields',

   methods: {
     getValue: function (field) {
       switch (field.data_type) {
         case 'DATE':
           return dateUtils.dateFormat(field.value)
         default:
           return field.value ? field.value : 'None'
       }
     }
   },

   props: {
     fields: {
       name: 'fields',
       default: []
     }
   }
 }
</script>

<style>
 .field-value {
   margin-top: 0.5rem;
   margin-bottom: 0.5rem;
 }

 @media(min-width: 992px) {
   .fields {
     margin-top: 1rem;
   }
 }
</style>
