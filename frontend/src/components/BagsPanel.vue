<template>
   <Panel header="Bag Details">
      <div class="wait"  v-if="submission.loadingBags" >
         <WaitSpinner message="Loading bag data..." />
      </div>
      <div v-else  class="bag-info">
          <div class="bag-summary" >
            <div>{{ submission.detail.bagCount }} bag(s) containing {{ submission.detail.fileCount }} file(s)</div>
            <div><b>Total size: {{ $formatFileSize(submission.detail.totalFileSize) }}</b></div>
         </div>
         <Tree :value="submission.bags">
            <template #bag="slotProps">
               <dl class="bag">
                  <dt>Bag:</dt>  
                  <dd>{{ slotProps.node.label }}</dd>
                  <dt>Files:</dt>  
                  <dd>{{ slotProps.node.children.length }}</dd>
                  <dt>Status:</dt>  
                  <dd>
                     <div class="status">
                        <span>{{ slotProps.node.data.status[0].status }}</span>
                     </div>
                  </dd>
                  <dt>Created:</dt>  
                  <dd>{{ $formatDateTime(slotProps.node.data.createdAt) }}</dd>
               </dl>
            </template>
            <template #file="slotProps">
               <dl class="file">
                  <dt>Filename:</dt>  
                  <dd>{{ slotProps.node.label }}</dd>
                  <dt>Size:</dt>  
                  <dd>{{ $formatFileSize(slotProps.node.data.fileSize) }}</dd>
                  <dt>Hash:</dt>  
                  <dd>{{ slotProps.node.data.hash }}</dd>
                  <dt>Created:</dt>  
                  <dd>{{ $formatDateTime(slotProps.node.data.createdAt) }}</dd>
               </dl>
            </template>
         </Tree>
      </div>
   </Panel>
</template>

<script setup>
import Panel from 'primevue/panel'
import { useSubmissionsStore } from "@/stores/submissions"
import Tree from 'primevue/tree'
import WaitSpinner from './WaitSpinner.vue'

const submission = useSubmissionsStore()
</script>

<style lang="scss" scoped>
.wait {
   text-align: center;
}
dl.file, dl.bag {
   padding: 0.5rem 1rem;
   border: 1px solid $uva-grey-100;
   border-radius: 0.3rem;
   grid-template-columns: max-content 1fr;
   display: inline-grid;
   grid-column-gap: 0.5rem;
   dd,dt {
      padding: 0.25rem 0;
    }
}
.bag-summary {
   display: flex;
   flex-flow: row wrap;
   justify-content: space-between;
}
.status {
   display: flex;
   flex-flow: row nowrap;
   justify-content: flex-start;
   align-items: center;
   span {
      text-transform: capitalize;
   }
}
dl.file {
   margin: 0 0 0 25px;
}
</style>