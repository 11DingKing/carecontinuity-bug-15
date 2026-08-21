package workercapacity

type Coordinator struct{ registry *TeamRegistry }

func NewCoordinator() *Coordinator {
	return &Coordinator{registry: NewTeamRegistry(AssignmentPolicy{LockScope: "task", HookEnabled: true})}
}
func (c *Coordinator) Assign(team, task string) bool    { return c.registry.Assign(team, task) }
func (c *Coordinator) Owner(team string) (string, bool) { return c.registry.Owner(team) }
func (c *Coordinator) SetBarrier(fn func())             { c.registry.BeforeStore = fn }
