package polymor

type coreOBJ interface {
	NewId()
	SetCreatedAt()
	SetUpdatedAt()
}
func SetDefult(o coreOBJ) {
	o.NewId()
	o.SetCreatedAt()
	o.SetUpdatedAt()
}