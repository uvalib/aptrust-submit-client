<template>
   <Panel header="Submission Approval">
      <div class="approve">
         <div class="storage">
            <label for="apt-storage">Storage</label>
            <Select inputId="apt-storage" v-model="storage" :options="storageOptions" placeholder="Select an option" fluid />
         </div>
         <div class="acts">
            <Button severity="danger" label="Cancel Submission" @click="cancelClicked"/>
            <Button label="Approve Submission" @click="approveClicked"/>
         </div>
      </div>
   </Panel>
</template>

<script setup>
import { ref,computed } from 'vue'
import Panel from 'primevue/panel'
import Select from 'primevue/select'
import { useSubmissionsStore } from "@/stores/submissions"
import { useSystemStore } from "@/stores/system"
import { useConfirm } from "primevue/useconfirm"

const system = useSystemStore()
const submission = useSubmissionsStore()
const confirm = useConfirm()

const storage = ref(submission.detail.storage)

const storageOptions = computed( () => {
   let out = [] 
   system.storageOptions.forEach( c=> {
      out.push(c.value)
   })
   return out
})

const cancelClicked = (()=> {
   confirm.require({
      message: "This APTrust submission will be canceled. Are you sure?",
      header: 'Confirm Cancel Submission',
      icon: 'pi pi-question-circle',
      rejectClass: 'p-button-secondary',
      accept: (  ) => {
        submission.cancel()
      },
   })
})

const approveClicked = (()=> {
   confirm.require({
      message: "All bags associated with this siubmission will be sent to APTrust. Are you sure?",
      header: 'Confirm APTrust Submission',
      icon: 'pi pi-question-circle',
      rejectClass: 'p-button-secondary',
      accept: (  ) => {
         submission.approve(storage.value)
      },
   })
})

</script>

<style lang="scss" scoped>
.approve {
   display: flex;
   flex-direction: column;
   gap: 20px;
}
.storage {
   display: flex;
   flex-flow: row nowrap;
   gap: 10px;
   justify-content: flex-start;
   align-items: center;
   label {
      font-weight: bold;
   }
}
.acts {
   display: flex;
   flex-flow: row wrap;
   justify-content: flex-end;
   gap: 10px;
}
</style>