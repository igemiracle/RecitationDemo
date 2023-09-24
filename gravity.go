package main

import (
	"math"
)

// let's place our gravity simulation functions here.
// simulateGravity takes as input an initial Universe object, a number of generations numGens, and a time parameter
func SimulateGravity(initialUniverse Universe, numGens int, time float64) []Universe {
	timePoints := make([]Universe, numGens+1)
	timePoints[0] = initialUniverse
	//range over the number of  generation and set the i-th Universe = updating the i-1 th Universe
	for i := 1; i < len(timePoints); i++ {
		timePoints[i] = UpdateUniverse(timePoints[i-1], time)
	}

	return timePoints
}

// updateUniverse returns a new Universe
func UpdateUniverse(currentUniverse Universe, time float64) Universe {
	//newUniverse := currentUniverse //BAD
	newUniverse := CopyUniverse(currentUniverse)

	//range over the bodies of the universe, and update their position/velocity/acceleration
	for i := range newUniverse.bodies {
		newUniverse.bodies[i].acceleration = UpdateAcceleration(currentUniverse, newUniverse.bodies[i])
		newUniverse.bodies[i].velocity = UpdateVelocity(newUniverse.bodies[i], time)
		newUniverse.bodies[i].position = UpdatePosition(newUniverse.bodies[i], time)
	}

	return newUniverse
}

// copyuniverse
func CopyUniverse(currentUniverse Universe) Universe {
	var newUniverse Universe

	newUniverse.width = currentUniverse.width

	//copy the body
	//first make a slice
	numBodies := len(currentUniverse.bodies)
	newUniverse.bodies = make([]Body, numBodies)

	//then copy every body's fields
	//range over all bodies
	for i := range newUniverse.bodies {
		newUniverse.bodies[i] = CopyBody(currentUniverse.bodies[i])
	}

	return newUniverse
}

// copybody
func CopyBody(oldBody Body) Body {
	var newBody Body
	newBody.name = oldBody.name
	newBody.mass = oldBody.mass
	newBody.radius = oldBody.radius
	newBody.red = oldBody.red
	newBody.green = oldBody.green
	newBody.blue = oldBody.blue

	newBody.position.x = oldBody.position.x
	newBody.position.y = oldBody.position.y
	newBody.velocity.x = oldBody.velocity.x
	newBody.velocity.y = oldBody.velocity.y
	newBody.acceleration.x = oldBody.acceleration.x
	newBody.acceleration.y = oldBody.acceleration.y

	return newBody
}

// takes as input a universe object and a body in that universe
// returns the net acceleration due to the force of gravity of the body computed
func UpdateAcceleration(currentUniverse Universe, b Body) OrderedPair {
	var accel OrderedPair

	force := ComputeNetForce(currentUniverse, b)
	//now, use Newton's law F=ma or a=F/m

	//split acceleration compenetwise
	accel.x = force.x / b.mass
	accel.y = force.y / b.mass

	return accel
}

func ComputeNetForce(currentUniverse Universe, b Body) OrderedPair {
	var netForce OrderedPair

	//range over all the bodies other than b, and pass computing the force of gravity to a subroutine, then add in compoenets to net force vector
	for i := range currentUniverse.bodies {
		if currentUniverse.bodies[i] != b {
			force := ComputeForce(b, currentUniverse.bodies[i])
			netForce.x += force.x
			netForce.y += force.y
		}
	}

	return netForce
}

// input: 2 Body objects b1 and b2
// return: an OrderedPair corresponding to the force of gravity of b2 acting on b1
func ComputeForce(b1, b2 Body) OrderedPair {
	var force OrderedPair

	//F=G*m1*m2/d^2
	dist := Distance(b1.position, b2.position)

	F := G * b1.mass * b2.mass / (dist * dist)

	//trick deltaX deltaY
	deltaX := b2.position.x - b1.position.x
	deltaY := b2.position.y - b1.position.y

	force.x = F * deltaX / dist
	force.y = F * deltaY / dist

	return force

}

func UpdateVelocity(b Body, time float64) OrderedPair {

	var vel OrderedPair

	vel.x = b.velocity.x + b.acceleration.x*time
	vel.y = b.velocity.y + b.acceleration.y*time

	return vel
}

func UpdatePosition(b Body, time float64) OrderedPair {

	var pos OrderedPair
	pos.x = b.position.x + b.velocity.x*time + 0.5*b.acceleration.x*time*time
	pos.y = b.position.y + b.velocity.y*time + 0.5*b.acceleration.y*time*time

	return pos

}

// Distance takes two position ordered pairs and it returns the distance between these two points in 2-D space.
func Distance(p1, p2 OrderedPair) float64 {
	// this is the distance formula from days of precalculus long ago ...
	deltaX := p1.x - p2.x
	deltaY := p1.y - p2.y
	return math.Sqrt(deltaX*deltaX + deltaY*deltaY)
}
