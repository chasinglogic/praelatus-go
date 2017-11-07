<template>
  <div class="notification" v-html="renderedNotification">
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
     renderedNotification: function () {
       let s = this.notification.body.split(' ')
       let username = '<router-link :to=\'/users/' + s[0] + '\'>' +
                      s[0] + '</router-link>'
       let rest = s.slice(1)
       return username + ' ' + rest.join(' ')
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
