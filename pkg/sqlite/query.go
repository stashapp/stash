package sqlite

type repositoryQueryBuilder struct {
	queryBuilder
	repository *repository
}

func (qb repositoryQueryBuilder) executeFind() ([]int, int, error) {
	return qb.repository.executeFindQuery(qb.body, qb.args, qb.sortAndPagination, qb.whereClauses, qb.havingClauses)
}
