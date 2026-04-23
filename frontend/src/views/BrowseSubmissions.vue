<template>
   <div class="dashboard">
      <h1>Browse APTrust Submissions</h1>   
      <WaitSpinner v-if="system.working" :overlay="true" message="<div>Please wait...</div><p>Searching APTrust submissions</p>" />
      <div class="content" v-else>
         CONTENT Here!
      </div>
   </div>
</template>

<script setup>
import WaitSpinner from "@/components/WaitSpinner.vue"
import { useSystemStore } from "@/stores/system"
const system = useSystemStore()

</script>

<style lang="scss">
.dashboard {
   margin: 0 auto 50px;
   min-height: 600px;
   text-align: left;
}
@media only screen and (min-width: 768px) {
   .dashboard {
       width: 90%;
   }
}
@media only screen and (max-width: 768px) {
   .dashboard {
      width: 95%;
   }
}
</style>