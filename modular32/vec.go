package modular32

import (
	mgl "github.com/go-gl/mathgl/mgl32"
)

// NewVec2Modulus creates a new 2d Vector Modulus
func NewVec2Modulus(vec mgl.Vec2) Vec2Modulus {
	return Vec2Modulus{
		x: NewModulus(vec[0]),
		y: NewModulus(vec[1]),
	}
}

// Vec2Modulus defines a modulus for 2d vectors
type Vec2Modulus struct {
	x Modulus
	y Modulus
}

// Congruent performs Congruent() on all axis
func (m Vec2Modulus) Congruent(vec mgl.Vec2) mgl.Vec2 {
	return mgl.Vec2{
		m.x.Congruent(vec[0]),
		m.y.Congruent(vec[1]),
	}
}

// Dist returns the distance and direction of v1 to v2
// It picks the shortest distance.
func (m Vec2Modulus) Dist(v1, v2 mgl.Vec2) mgl.Vec2 {
	return mgl.Vec2{
		m.x.Dist(v1[0], v2[0]),
		m.y.Dist(v1[1], v2[1]),
	}
}

// GetCongruent returns the vector closest to v1 that is congruent to v2
func (m Vec2Modulus) GetCongruent(v1, v2 mgl.Vec2) mgl.Vec2 {
	return mgl.Vec2{
		m.x.GetCongruent(v1[0], v2[0]),
		m.y.GetCongruent(v1[1], v2[1]),
	}
}

// NewVec3Modulus creates a new 3d Vector Modulus
func NewVec3Modulus(vec mgl.Vec3) Vec3Modulus {
	return Vec3Modulus{
		x: NewModulus(vec[0]),
		y: NewModulus(vec[1]),
		z: NewModulus(vec[2]),
	}
}

// Vec3Modulus defines a modulus for 3d vectors
type Vec3Modulus struct {
	x Modulus
	y Modulus
	z Modulus
}

// Congruent performs Congruent() on all axis
func (m Vec3Modulus) Congruent(vec mgl.Vec3) mgl.Vec3 {
	return mgl.Vec3{
		m.x.Congruent(vec[0]),
		m.y.Congruent(vec[1]),
		m.z.Congruent(vec[2]),
	}
}

// Dist returns the distance and direction of v1 to v2
// It picks the shortest distance.
func (m Vec3Modulus) Dist(v1, v2 mgl.Vec3) mgl.Vec3 {
	return mgl.Vec3{
		m.x.Dist(v1[0], v2[0]),
		m.y.Dist(v1[1], v2[1]),
		m.z.Dist(v1[2], v2[2]),
	}
}

// GetCongruent returns the vector closest to v1 that is congruent to v2
func (m Vec3Modulus) GetCongruent(v1, v2 mgl.Vec3) mgl.Vec3 {
	return mgl.Vec3{
		m.x.GetCongruent(v1[0], v2[0]),
		m.y.GetCongruent(v1[1], v2[1]),
		m.z.GetCongruent(v1[2], v2[2]),
	}
}

// NewVec4Modulus creates a new 4d Vector Modulus
func NewVec4Modulus(vec mgl.Vec4) Vec4Modulus {
	return Vec4Modulus{
		x: NewModulus(vec[0]),
		y: NewModulus(vec[1]),
		z: NewModulus(vec[2]),
		w: NewModulus(vec[3]),
	}
}

// Vec4Modulus defines a modulus for 4d vectors
type Vec4Modulus struct {
	x Modulus
	y Modulus
	z Modulus
	w Modulus
}

// Congruent performs Congruent() on all axis
func (m Vec4Modulus) Congruent(vec mgl.Vec4) mgl.Vec4 {
	return mgl.Vec4{
		m.x.Congruent(vec[0]),
		m.y.Congruent(vec[1]),
		m.z.Congruent(vec[2]),
		m.w.Congruent(vec[3]),
	}
}

// Dist returns the distance and direction of v1 to v2
// It picks the shortest distance.
func (m Vec4Modulus) Dist(v1, v2 mgl.Vec4) mgl.Vec4 {
	return mgl.Vec4{
		m.x.Dist(v1[0], v2[0]),
		m.y.Dist(v1[1], v2[1]),
		m.z.Dist(v1[2], v2[2]),
		m.z.Dist(v1[3], v2[3]),
	}
}

// GetCongruent returns the vector closest to v1 that is congruent to v2
func (m Vec4Modulus) GetCongruent(v1, v2 mgl.Vec4) mgl.Vec4 {
	return mgl.Vec4{
		m.x.GetCongruent(v1[0], v2[0]),
		m.y.GetCongruent(v1[1], v2[1]),
		m.z.GetCongruent(v1[2], v2[2]),
		m.w.GetCongruent(v1[3], v2[3]),
	}
}
