<template>
  <div class="fields">
    <template v-for="field in ticket.fields" >
      <p>
        <span class="bold">{{ field.name }}</span>
        {{ field.value }}
      </p>
    </template>
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
