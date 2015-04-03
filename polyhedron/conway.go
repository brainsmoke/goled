package polyhedron

var ConwayModels = map[string] func() []Face {

	"T": TetrahedronFaces,
	"O": OctahedronFaces,
	"C": CubeFaces,
	"D": DodecahedronFaces,
	"I": IcosahedronFaces,
	"tT": TruncatedTetrahedronFaces,
	"tC": TruncatedCubeFaces,
	"bC": TruncatedCuboctahedronFaces,
	"tO": TruncatedOctahedronFaces,
	"tD": TruncatedDodecahedronFaces,
	"bD": TruncatedIcosidodecahedronFaces,
	"tI": TruncatedIcosahedronFaces,
	"aC": CuboctahedronFaces,
	"aD": IcosidodecahedronFaces,
	"eC": RhombicuboctahedronFaces,
	"eD": RhombicosidodecahedronFaces,
	"sC": SnubCubeFaces,
	"sD": SnubDodecahedronFaces,
	"kT": TriakisTetrahedronFaces,
	"kO": TriakisOctahedronFaces,
	"mC": DisdyakisDodecahedronFaces,
	"kC": TetrakisHexahedronFaces,
	"kI": TriakisIcosahedronFaces,
	"mD": DisdyakisTriacontahedronFaces,
	"kD": PentakisDodecahedronFaces,
	"jC": RhombicDodecahedronFaces,
	"jD": RhombicTriacontahedronFaces,
	"oC": DeltoidalIcositetrahedronFaces,
	"oD": DeltoidalHexecontahedronFaces,
	"gC": PentagonalIcositetrahedronFaces,
	"gD": PentagonalHexecontahedronFaces,

}

