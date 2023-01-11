package entry

// CheckCollision 检查子弹和外星人之间是否有碰撞,将两个实体看做两个矩形，计算判断是否处于碰撞
func CheckCollision(entityA, entityB Entity) bool {
	rec1 := []float64{entityA.X(), entityA.Y() + float64(entityA.Height()), entityA.X() + float64(entityA.Width()), entityA.Y()}
	rec2 := []float64{entityB.X(), entityB.Y() + float64(entityB.Height()), entityB.X() + float64(entityB.Width()), entityB.Y()}

	x_overlap := !(rec1[2] <= rec2[0] || rec2[2] <= rec1[0])
	y_overlap := !(rec1[3] >= rec2[1] || rec2[3] >= rec1[1])
	return x_overlap && y_overlap
}
