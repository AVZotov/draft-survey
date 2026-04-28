function resultsPage(alpine, lastIndex) {
  return {
    data: JSON.parse(alpine),
    firstIndex: 0,
    secondIndex: lastIndex,
    drafts: [],
    results: [],

    init() {
      this.drafts = this.data.drafts.map((draft, i) => ({
        draft,
        index: i,
        label: this.labelFor(draft, i)
      }))
      this.results = this.data.draft_totals

      this.$watch('firstIndex', value => {
        this.secondIndex = value + 1
      })
    },

    labelFor(draft, index) {
      if (draft.type === 'initial') return 'Initial'
      if (draft.type === 'final') return 'Final'
      return 'Intermediate ' + index
    },

    get firstSelector() {
      return this.drafts.slice(0, this.drafts.length - 1)
    },

    get secondSelector() {
      return this.drafts.filter(item => item.index > this.firstIndex)
    },

    openSections: [false, true, false, false, false, false, false, false, false, false],

    typeClass(type) {
      const map = { initial: 'draft-col--ini', final: 'draft-col--fin', intermediate: 'draft-col--mid' }
      return map[type] || 'draft-col--ini'
    },
  }
}