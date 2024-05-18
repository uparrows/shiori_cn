var template = `
<div class="pagination-box">
	<p>页面</p>
	<input type="text" 
		placeholder="1" 
		:value="page" 
		@focus="$event.target.select()" 
		@keyup.enter="changePage($event.target.value)" 
		:disabled="editMode">
	<p>{{maxPage}}</p>
	<div class="spacer"></div>
	<template v-if="!editMode">
		<a v-if="page > 2" title="转到第一页" @click="changePage(1)">
			<i class="fas fa-fw fa-angle-double-left"></i>
		</a>
		<a v-if="page > 1" title="转到上一页" @click="changePage(page-1)">
			<i class="fa fa-fw fa-angle-left"></i>
		</a>
		<a v-if="page < maxPage" title="转到下一页" @click="changePage(page+1)">
			<i class="fa fa-fw fa-angle-right"></i>
		</a>
		<a v-if="page < maxPage - 1" title="转到最后一页" @click="changePage(maxPage)">
			<i class="fas fa-fw fa-angle-double-right"></i>
		</a>
	</template>
</div>`;

export default {
	template: template,
	props: {
		page: Number,
		maxPage: Number,
		editMode: Boolean,
	},
	methods: {
		changePage(page) {
			page = parseInt(page, 10) || 0;
			if (page >= this.maxPage) page = this.maxPage;
			else if (page <= 1) page = 1;

			this.$emit("change", page);
		},
	},
};
