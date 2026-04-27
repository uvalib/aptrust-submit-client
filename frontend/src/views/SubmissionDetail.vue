<template>
   <div class="details">
      <h1>Submission {{ route.params.id }}</h1>  
   </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useSubmissionsStore } from "@/stores/submissions"
import { useSystemStore } from "@/stores/system"
import { useRoute } from 'vue-router'

const submissionStore = useSubmissionsStore()
const system = useSystemStore()
const route = useRoute()

onMounted( () => {
   submissionStore.getSubmissionDetail( route.params.id )
})
</script>

<style lang="scss">
.details {
   margin: 0 auto 50px;
   min-height: 600px;
   text-align: left;
}
@media only screen and (min-width: 768px) {
   .details {
       width: 90%;
   }
}
@media only screen and (max-width: 768px) {
   .details {
      width: 95%;
   }
}
</style>
   