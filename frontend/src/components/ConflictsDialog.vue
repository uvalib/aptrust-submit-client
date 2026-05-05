<template>
   <Button size="small" severity="secondary" @click="showDialog = true" label="Details" iconPos="right"/>
   <Dialog v-model:visible="showDialog" :modal="true" header="Submission Conflicts" style="width: 95%;">
      <DataTable :value="submission.detail.conflicts" dataKey="id" v-model:expandedRows="expandedRows"
         stripedRows showGridlines responsiveLayout="scroll" :alwaysShowPaginator="false"
         :lazy="false" :paginator="true" :rows="30" :totalRecords="submission.detail.conflicts.length"
         paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
         currentPageReportTemplate="{first} - {last} of {totalRecords}" paginatorPosition="both"
      >
         <Column :expander="true" headerStyle="width: 3rem" />
         <Column field="createdAt" header="Timestamp">
            <template #body="slotProps">
               {{ $formatDateTime(slotProps.data.createdAt) }}
            </template>
         </Column>
         <Column header="Filename" >
            <template #body="slotProps">
               <span>{{ slotProps.data.newFile.fileName }}</span>
            </template>
         </Column>
         <Column field="basis" header="Type" />
         <Column header="Conflicting Filename" >
            <template #body="slotProps">
               <span v-if="conflictFileExists(slotProps.data)">{{ conflictFile(slotProps.data).fileName }}</span>
               <span v-else>Not Found</span>
            </template>
         </Column>
         <Column field="ignored" header="Ignored" class="centered">
            <template #body="slotProps">
               <span class="check pi pi-circle-fill" :class="{ignored: slotProps.data.ignored}"></span>
            </template>
         </Column>
         <template #expansion="slotProps">
            <div class="file-details">
               <Panel header="Submission File">
                  <dl>
                     <dt>Filename</dt>  
                     <dd>{{ slotProps.data.newFile.fileName }}</dd>
                     <dt>Size</dt>  
                     <dd>{{ $formatFileSize(slotProps.data.newFile.fileSize) }}</dd>
                     <dt>Hash</dt>  
                     <dd>{{ slotProps.data.newFile.hash }}</dd>
                     <dt>Bag name</dt>  
                     <dd>{{ slotProps.data.newFile.bagName }}</dd>
                  </dl>
               </Panel>
               <Panel header="Conflicting File">
                  <dl v-if="conflictFileExists(slotProps.data)">
                     <dt>Filename</dt>  
                     <dd>{{ conflictFile(slotProps.data).fileName }}</dd>
                     <dt>Size</dt>  
                     <dd>{{ $formatFileSize( conflictFile(slotProps.data).fileSize ) }}</dd>
                     <dt>Hash</dt>  
                     <dd>{{ conflictFile(slotProps.data).hash }}</dd>
                     <dt>Bag name</dt>  
                     <dd>{{ conflictFile(slotProps.data).bagName }}</dd>
                  </dl>
                  <div v-else>Not Found</div>
               </Panel>
            </div>
         </template>
      </DataTable>       
   </Dialog>
</template>

<script setup>
import Dialog from 'primevue/dialog'
import { ref } from 'vue'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Panel from 'primevue/panel'
import { useSubmissionsStore } from "@/stores/submissions"

const submission = useSubmissionsStore()
const showDialog = ref(false)
const expandedRows = ref([])

const conflictFileExists = ( (conflict) => {
   return (conflict.aptConflict || conflict.localConflict)
})

const conflictFile = ( (conflict) => {
   if (conflict.aptConflict) return conflict.aptConflict
   return conflict.localConflict
})

</script>

<style lang="scss" scoped>
.file-details {
   display: flex;
   flex-flow: row wrap;
   gap: 15px;
   .p-panel {
      flex-grow: 1;
   }
   dt {
      font-weight: bold;
   }
   dd {
      margin-bottom: 15px;
   }
}
:deep(.check) {
   color: $uva-red-A;
   font-size: 1.3rem;
}
:deep(.check.ignored) {
   color: $uva-green;
   font-size: 1.3rem;
}
:deep(   .centered) {
   text-align: center;
}

</style>
