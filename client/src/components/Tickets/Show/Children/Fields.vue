<template>
  <div class="card fields">
    <h2 class="card-header" >
      Fields
    </h2>

    <div class="card-block" >
      <div v-for="field in fields">
        <p class="bold">{{ field.name }}</p>
        <p>{{ getValue(field) }}</p>
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
           return field.value
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
 .fields {
   margin-top: 1rem;
 }
</style>
