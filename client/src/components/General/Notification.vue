<template>
  <div class="notification card" >
    <div class="card-block">
      <user-stub :username="notification.actioningUser" />
      {{ actionText }}
      <router-link :to="'/tickets/' + notification.actionedTicket">
        {{ notification.actionedTicket }}
      </router-link>
    </div>
  </div>
</template>

<script>
 import UserStub from '@/components/Users/Stub'

 export default {
   name: 'notification',

   components: {
     UserStub
   },

   computed: {
     actionText: function () {
       return this.notification.body
                  .replace(this.notification.actioningUser, '')
                  .replace(this.notification.actionedTicket, '')
     }
   },

   props: {
     notification: {
       type: Object,
       required: true
     }
   }
 }
</script>

<style>
 .notification .avatar {
   height: 20px;
   width: 20px;
 }

 .notification {
   margin-bottom: 0.5rem;
 }

 .notification .card-block {
   padding: 0.5rem;
 }
</style>
