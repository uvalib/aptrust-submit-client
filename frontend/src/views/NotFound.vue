<template>
   <h1>Not Found</h1>
   <div class="column gap-4 items-center text-xl">
      <div>Sorry, the URL you requested could not be found!</div>
      <UButton label="Return to Home Page"  @click="homeClicked()" />
   </div>
</template>

<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()
function homeClicked() {
   router.push("/")
}
</script>

<style scoped lang="scss">
</style>
