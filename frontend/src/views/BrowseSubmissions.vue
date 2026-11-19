<template>
   <div class="dashboard">
      <h1>Browse AP Trust Submissions</h1>   
       <WaitSpinner v-if="submissionStore.working" message="Loading submissions..." :overlay="true" />
      <div class="row-between p-2 bg-white">
         <UCheckbox v-model="submissionStore.includeAutoApproved" label="Include auto-approved submissions"
            size="xl" @update:model-value="includeAutoApproveChanged()"/>
         <div class="row-right gap-4">
            <UPagination color="neutral" variant="ghost" 
               v-model:page="submissionStore.currPage" :items-per-page="submissionStore.pageSize" 
               :total="submissionStore.total" @update:page="pageChanged"
            />
            <USelect v-model="submissionStore.pageSize" :items="[30,50,100]" @change="perPageChanged" />
            <span>Total: {{ submissionStore.total }}</span>
         </div>
         <div class="row-right gap-2">
            <UInput v-model="submissionStore.query" placeholder="Search submissions..." @keypress="searchKeyPressed($event)" />
            <UButton color="secondary" label="Reset" @click="resetSearchClicked()"/>
         </div>
      </div>
      <UTable :data="submissionStore.searchHits" :columns="columns" v-model:sorting="sorting" @select="onSelect" sticky class="h-[calc(100vh-185px)]" >
         <template #status-header="{ column }">
            <FilterPopover :label="column.columnDef.header" :applied="isFilterApplied(column.id)" @clear="clearFilter(column.id)" @apply="applyFilters">
               <USelectMenu v-model="filters.status.value" :items="system.submissionStatuses" class="w-full" 
                  placeholder="Select status..."/>   
            </FilterPopover>
         </template>
         <template #status-cell="{ row }">
            <div :class="statusClass(row.original.status)" v-if="row.original.status">
               <span>{{ row.original.status.replace('-', ' ') }}</span>
               <UIcon v-if="row.original.status=='pending-approval'" name="i-lucide-loader-circle" class="size-5 animate-spin" size="20" />
               <UIcon v-if="row.original.status=='error'" name="i-lucide-circle-letter-x" size="20" />
               <UIcon v-if="row.original.status=='abandoned'" name="i-lucide-circle-check" size="20" />
            </div>
            <span v-else class="none">Unknown</span>   
         </template>
         <template #client-header="{ column }">
            <FilterPopover :label="column.columnDef.header" :applied="isFilterApplied(column.id)" @clear="clearFilter(column.id)" @apply="applyFilters">
               <USelectMenu v-model="filters.client.value" :items="clients" class="w-full" 
                  placeholder="Select client..."/>   
            </FilterPopover>
         </template>
         <template #createdAt-header="{ column }">
            <UButton color="neutral" variant="ghost" :icon="sortIcon(column)" @click="sortClicked(column)"> 
               <span class="text-start">{{ column.columnDef.header }}</span>
            </UButton>
         </template>
         <template #createdAt-cell="{ row }">
            {{ $formatDateTime(row.original.createdAt) }}
         </template>
         <template #statusCreatedAt-header="{ column }">
            <UButton color="neutral" variant="ghost" :icon="sortIcon(column)" @click="sortClicked(column)"> 
               <span class="text-start">{{ column.columnDef.header }}</span>
            </UButton>
         </template>
         <template #statusCreatedAt-cell="{ row }">
            {{ $formatDateTime(row.original.statusCreatedAt) }}
         </template>
      </UTable>
   </div>
</template>

<script setup>
import { onMounted, ref, computed, onBeforeMount } from 'vue'
import { useSubmissionsStore } from "@/stores/submissions"
import { useSystemStore } from "@/stores/system"
import { useRouter, useRoute } from 'vue-router'

const submissionStore = useSubmissionsStore()
const system = useSystemStore()
const router = useRouter()
const route = useRoute()

const sorting = ref([ { id: 'createdAt', desc: true } ])

const columns = [
   { accessorKey: 'status', header: 'Status', meta: { class: { th: 'w-[200px]!', td: 'w-[200px]!' } } }, 
   { accessorKey: 'collectionName', header: 'Submission' },
   { accessorKey: 'client', header: 'Source' },
   { accessorKey: 'createdAt', header: 'Created' },
   { accessorKey: 'statusCreatedAt', header: 'Updated' },
]

const filters = ref({
   status: {value: null, mode: "equals"},
   client: {value: null, mode: "equals"},
})

const clients = computed( () => {
   let out = [] 
   system.clients.forEach( c=> {
      out.push(c.name)
   })
   return out
})

onBeforeMount( () => {
   // pull query params and update store filter, query and pagination settings
   // they will be used in onMounted when submissions are requested
   if ( route.query.autoapproved ) {
      submissionStore.includeAutoApproved = true 
   }
   if ( route.query.q ) {
      submissionStore.query = route.query.q 
   }
   if ( route.query.offset ) {
      submissionStore.offset = parseInt(route.query.offset, 10)
   }
   if ( route.query.limit ) {
      submissionStore.pageSize = parseInt(route.query.limit, 10)
   }
   if ( route.query.sort ) {
      submissionStore.sortField = route.query.sort
   }
   if ( route.query.order ) {
      submissionStore.sortOrder = route.query.order
   }
   if ( route.query.filter ) { 
      let filterJson = JSON.parse(route.query.filter)
      filterJson.forEach( fj => {
         let bits = fj.split("=")
         submissionStore.filters.push( {field: bits[0], match: "equals", value: bits[1]})
         if (bits[0] == "client") {
            filters.value.client.value = bits[1]
         } else if (bits[0] == "status") {
            filters.value.clstatusient.value = bits[1]
         }
      })
   }
})

onMounted( () => {
   submissionStore.getSubmissions()
})

const resetSearchClicked = (() => {
   let query = Object.assign({}, route.query)
   delete query.autoapproved
   delete query.q
   delete query.offset
   delete query.pageSize
   delete query.sortField
   delete query.sortOrder
   delete query.filter
   router.push({query})

   // clear the primevue datatables filter model
   for (let key in filters.value) {
      filters.value[key].value = null
   }
   
   submissionStore.$reset()
   submissionStore.getSubmissions()
})

const searchToQueryParams = (() => {
   let query = Object.assign({}, route.query)
   delete query.autoapproved
   delete query.q
   delete query.offset
   delete query.pageSize
   delete query.sortField
   delete query.sortOrder
   delete query.filter
   if ( submissionStore.includeAutoApproved ) {
      query.autoapproved = 1
   }
   if ( submissionStore.query != "" ) {
      query.q = submissionStore.query
   }
   if (submissionStore.offset > 0 ) {
      query.offset = submissionStore.offset
   }
   query.limit = submissionStore.pageSize
   query.sort = submissionStore.sortField
   query.order = submissionStore.sortOrder
   if (submissionStore.filters.length > 0) {
      query.filter = submissionStore.filtersAsQueryParam
   }
   router.push({query})
   submissionStore.getSubmissions()
})

const includeAutoApproveChanged = (() => {
   searchToQueryParams()
})

const statusClass = ( (status)=> {
   if (status == "pending-approval") return "status pending"
   if (status == "complete" || status == "abandoned") return "status done"
   if (status == "error") return "status error"
   return  "status"
})

const onSelect = ((_event, row) => {
   const subID = row.original.identifier
   router.push(`/submissions/${subID}`)
})

const sortIcon = ((column) => {
   const isSorted = column.getIsSorted()
   if (isSorted) {
      if (isSorted === 'asc') return 'i-lucide-arrow-up-narrow-wide'
      return 'i-lucide-arrow-down-wide-narrow'
   }
   return 'i-lucide-arrow-up-down'
})
const sortClicked = (async (column) => {
   const isSorted = column.getIsSorted()
   submissionStore.sortField = column.id
   submissionStore.sortOrder = "asc"
   if (isSorted === "asc") {
      submissionStore.sortOrder = "desc"   
   }
   searchToQueryParams()
   column.toggleSorting(column.getIsSorted() === 'asc')
})

const searchKeyPressed = ((event) => {
   submissionStore.offset = 0
   if (event.keyCode == 13) {
      searchToQueryParams()
   }
})

const isFilterApplied = ((name) => {
   return filters.value[name].value != null 
})
const clearFilter = ((name) => {
   filters.value[name].value = null 
   submissionStore.resetSearch()
   searchToQueryParams()
})
const applyFilters = (() => {
   submissionStore.resetSearch()
   Object.entries(filters.value).forEach(([key, data]) => {
      if (data.value && data.value != "") {
         submissionStore.filters.push({field: key, match: data.mode, value: data.value})
      }
   })
   searchToQueryParams() 
})

const perPageChanged = (() => {
   submissionStore.currPage = 1
   pageChanged()
})
const pageChanged = (() => {
   submissionStore.offset = (submissionStore.currPage-1) * submissionStore.pageSize
   searchToQueryParams()
})

</script>

<style lang="scss">
.dashboard {
   margin: 0 auto 50px;
   min-height: 600px;
   text-align: left;
   .none {
      color: $uva-grey-A;
      font-style: italic;
   }
   .status {
      text-transform: capitalize;
      display: flex;
      flex-flow: row nowrap;
      justify-content: space-between;
      align-items: center;
      gap: 0.5rem;
      padding: 0.25rem .75rem; 
      background-color: $uva-grey-200;
      border: 1px solid $uva-grey;
      border-radius: 1rem;
   }
   .status.error {
       background-color: $uva-red-A;
       border-color: $uva-red-B;
       color: white;
       font-weight: bold;
   }
   .status.pending {
       background-color: $uva-yellow-100;
       border-color: $uva-yellow;
   }
   .status.done {
       background-color: $uva-green-A;
       border-color: $uva-green-B;
       color: white;
       font-weight: bold;
   }
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