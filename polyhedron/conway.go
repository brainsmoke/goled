package polyhedron

var ConwayModels = map[string] func() Solid {

	"T": Tetrahedron,
	"O": Octahedron,
	"C": Cube,
	"D": Dodecahedron,
	"I": Icosahedron,
	"tT": TruncatedTetrahedron,
	"tC": TruncatedCube,
	"bC": TruncatedCuboctahedron,
	"tO": TruncatedOctahedron,
	"tD": TruncatedDodecahedron,
	"bD": TruncatedIcosidodecahedron,
	"tI": TruncatedIcosahedron,
	"aC": Cuboctahedron,
	"aD": Icosidodecahedron,
	"eC": Rhombicuboctahedron,
	"eD": Rhombicosidodecahedron,
	"sC": SnubCube,
	"sD": SnubDodecahedron,
	"kT": TriakisTetrahedron,
	"kO": TriakisOctahedron,
	"mC": DisdyakisDodecahedron,
	"kC": TetrakisHexahedron,
	"kI": TriakisIcosahedron,
	"mD": DisdyakisTriacontahedron,
	"kD": PentakisDodecahedron,
	"jC": RhombicDodecahedron,
	"jD": RhombicTriacontahedron,
	"oC": DeltoidalIcositetrahedron,
	"oD": DeltoidalHexecontahedron,
	"gC": PentagonalIcositetrahedron,
	"gD": PentagonalHexecontahedron,

}

