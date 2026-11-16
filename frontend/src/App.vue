<template>
   <UApp :toaster="toaster">
      <UHeader mode="slideover" id="uva-header" title="UVA Library" to="https://library.virginia.edu">
         <template #title>
            <div class="library-link">
               <UvaLibraryLogo />
            </div>
         </template>
         <template #right>
            <div class="site-link">
               <RouterLink to="/">AP Trust Submission</RouterLink>
            </div>
         </template>

         <template #bottom>
            <div class="p-2 row-right gap-2 bg-brand-grey-200 border-b-brand-grey-100" v-if="system.isSignedIn">
               <span>Signed in as {{ system.signedInUser }}</span>
               <UButton icon="i-lucide-log-out" color="secondary" @click="system.signOut()"/> 
            </div> 
         </template>
      </UHeader>

      <UMain>
         <RouterView  v-if="configuring == false"/>
         <span id="new-window" tabindex="-1" class="screen-reader-text">opens in a new window</span>
      </UMain>

      <LibraryFooter />
   </UApp>

   <UModal v-model:open="system.showError" :modal="true" :dismissible="false" title="System Error">
      <template #body>
         <div style="text-align: left" v-html="system.error"></div>
      </template>
   </UModal>
</template>

<script setup>
import { onBeforeMount, ref, watch } from 'vue'
import UvaLibraryLogo from "@/components/UvaLibraryLogo.vue"
import LibraryFooter from "@/components/LibraryFooter.vue"
import { RouterView } from 'vue-router'
import { useSystemStore } from "@/stores/system"

const toast = useToast()
const toaster = { duration: 5000, position: "top-center" }

const system = useSystemStore()

const configuring = ref(true)

watch(() => system.toast.show, (newShow) => {
   if ( newShow == true) {
      let color = "primary"
      if ( system.toast.error) {
        color = "error"
      }
      console.log("ADD TOAST")
      toast.add({
         title: system.toast.summary,
         description: system.toast.message,
         color: color
      })
      system.clearToastMessage()
   }
})

onBeforeMount( async () => {
   await system.getConfig()
   configuring.value = false
})

const errorClosed = (() => {
   system.setError("")
   system.showError = false
})
</script>

<style lang="scss">
div.library-link {
   width: 220px;
}
div.site-link {
   text-align: right;
   font-size: 1.5em;
   font-weight: bold;
   a {
      text-decoration: none;
      color: white;
      &:hover {
         text-decoration: underline;
      }
   }
}
</style>
