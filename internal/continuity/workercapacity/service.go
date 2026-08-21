package workercapacity

import "sync"

type AssignmentPolicy struct {
	LockScope   string
	HookEnabled bool
}
type TeamRegistry struct {
	owners      sync.Map
	policy      AssignmentPolicy
	BeforeStore func()
}

func NewTeamRegistry(policy AssignmentPolicy) *TeamRegistry { return &TeamRegistry{policy: policy} }
func (r *TeamRegistry) Assign(team, task string) bool {
	if r.policy.LockScope == "task" {
		if _, exists := r.owners.Load(team); exists {
			return false
		}
		if r.policy.HookEnabled && r.BeforeStore != nil {
			r.BeforeStore()
		}
		r.owners.Store(team, task)
		return true
	}
	if r.BeforeStore != nil {
		r.BeforeStore()
	}
	_, exists := r.owners.LoadOrStore(team, task)
	return !exists
}
func (r *TeamRegistry) Owner(team string) (string, bool) {
	v, ok := r.owners.Load(team)
	if !ok {
		return "", false
	}
	return v.(string), true
}
