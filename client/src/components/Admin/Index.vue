<!-- Copyright 2017 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
     Use of this source code is governed by the AGPLv3 license that can be found in
     the LICENSE file. -->

<template>
  <div>
    <sidebar>
      <h5>System Administration</h5>
      <sidebar-links>
        <router-link to="/admin/workflows">
          Workflows
        </router-link>
      </sidebar-links>
    </sidebar>
  </div>
</template>

<script>
 import Sidebar from '@/components/General/Sidebar'
 import SidebarLinks from '@/components/General/SidebarLinks'

 export default {
   name: 'dashboard',
   components: {
     SidebarLinks,
     Sidebar
   }
 }
</script>
