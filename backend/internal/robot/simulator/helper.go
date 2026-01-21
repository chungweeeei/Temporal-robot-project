package simulator

import "math"

func transferOrientationToQuaternion(orientation float64) (qx, qy, qz, qw float64) {
	// Assuming orientation is in degrees and represents rotation around Z-axis (Yaw)
	// Convert degrees to radians
	radians := orientation * (math.Pi / 180.0)

	qx = 0.0
	qy = 0.0
	qz = math.Sin(radians / 2)
	qw = math.Cos(radians / 2)

	return qx, qy, qz, qw
}

func transferQuaternionToOrientation(qx, qy, qz, qw float64) float64 {
	// Assuming rotation around Z-axis only
	// Extract Yaw (Z-axis rotation) from quaternion
	// formula: yaw = atan2(2(wz + xy), 1 - 2(y^2 + z^2))
	// simpler for pure Z rotation: 2 * atan2(z, w)
	radians := 2 * math.Atan2(qz, qw)

	// Convert radians back to degrees
	degrees := radians * (180.0 / math.Pi)

	return degrees
}
