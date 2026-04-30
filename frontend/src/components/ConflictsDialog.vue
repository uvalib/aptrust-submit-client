<template>
   <Button size="small" severity="secondary" @click="showDialog = true" label="Details" iconPos="right"/>
   <Dialog v-model:visible="showDialog" :modal="true" header="Submission Conflicts" style="width: 95%;">
      <DataTable :value="props.conflicts" dataKey="id" v-model:expandedRows="expandedRows"
         stripedRows showGridlines responsiveLayout="scroll"
         :lazy="false" :paginator="true" :rows="30" :totalRecords="props.conflicts.length"
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
               <span>{{ conflictFile(slotProps.data).fileName }}</span>
            </template>
         </Column>
         <template #expansion="slotProps">
            <div class="file-details">
               <Panel header="Submission File">
                  <dl>
                     <dt>Filename</dt>  
                     <dd>{{ slotProps.data.newFile.fileName }}</dd>
                     <dt>Size</dt>  
                     <dd>{{ slotProps.data.newFile.fileSize }}</dd>
                     <dt>Hash</dt>  
                     <dd>{{ slotProps.data.newFile.hash }}</dd>
                     <dt>Bag name</dt>  
                     <dd>{{ slotProps.data.newFile.bagName }}</dd>
                  </dl>
               </Panel>
               <Panel header="Conflicting File">
                  <dl>
                     <dt>Filename</dt>  
                     <dd>{{ conflictFile(slotProps.data).fileName }}</dd>
                     <dt>Size</dt>  
                     <dd>{{ conflictFile(slotProps.data).fileSize }}</dd>
                     <dt>Hash</dt>  
                     <dd>{{ conflictFile(slotProps.data).hash }}</dd>
                     <dt>Bag name</dt>  
                     <dd>{{ conflictFile(slotProps.data).bagName }}</dd>
                  </dl>
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

const showDialog = ref(false)
const expandedRows = ref([])

const props = defineProps({
   conflicts: {
      type: Array,
      required: true
   },
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

</style>
