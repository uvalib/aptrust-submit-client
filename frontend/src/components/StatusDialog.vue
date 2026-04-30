<template>
   <Button size="small" icon="pi pi-history" severity="secondary" @click="showDialog = true" label="History" iconPos="right"/>
   <Dialog v-model:visible="showDialog" :modal="true" header="Submission Status History">
      <DataTable :value="submission.detail.status" stripedRows showGridlines responsiveLayout="scroll" 
         :alwaysShowPaginator="false"
         :lazy="false" :paginator="true" :rows="30" :totalRecords="submission.detail.status.length"
         paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
         currentPageReportTemplate="{first} - {last} of {totalRecords}" paginatorPosition="both"
      >
         <Column field="createdAt" header="Timestamp">
            <template #body="slotProps">
               {{ $formatDateTime(slotProps.data.createdAt) }}
            </template>
         </Column>
         <Column field="status" header="Status" >
            <template #body="slotProps">
               <span class="status">{{ slotProps.data.status.replace("-", " ") }}</span>
            </template>
         </Column>
      </DataTable>       
   </Dialog>
</template>

<script setup>
import Dialog from 'primevue/dialog'
import { ref } from 'vue'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import { useSubmissionsStore } from "@/stores/submissions"

const submission = useSubmissionsStore()
const showDialog = ref(false)
</script>

<style lang="scss" scoped>
.status {
   text-transform: capitalize;
}
</style>
