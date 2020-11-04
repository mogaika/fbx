package bfbx73

import (
	"github.com/mogaika/fbx"
)

func Author(author string) *fbx.Node                  { return fbx.NewNode("Author", author) }
func Colors(colors []float64) *fbx.Node               { return fbx.NewNode("Colors", colors) }
func Comment(comment string) *fbx.Node                { return fbx.NewNode("Comment", comment) }
func Connections() *fbx.Node                          { return fbx.NewNode("Connections") }
func Content(data []byte) *fbx.Node                   { return fbx.NewNode("Content", data) }
func Count(count int32) *fbx.Node                     { return fbx.NewNode("Count", count) }
func CreationTime(time string) *fbx.Node              { return fbx.NewNode("CreationTime", time) }
func CreationTimeStamp() *fbx.Node                    { return fbx.NewNode("CreationTimeStamp") }
func Creator(creator string) *fbx.Node                { return fbx.NewNode("Creator", creator) }
func Cropping(a1, a2, a3, a4 int32) *fbx.Node         { return fbx.NewNode("Cropping", a1, a2, a3, a4) }
func Culling(culling string) *fbx.Node                { return fbx.NewNode("Culling", culling) }
func Current(a1 string) *fbx.Node                     { return fbx.NewNode("Current", a1) }
func Day(day int32) *fbx.Node                         { return fbx.NewNode("Day", day) }
func Definitions() *fbx.Node                          { return fbx.NewNode("Definitions") }
func Deformer(id int64, a2, a3 string) *fbx.Node      { return fbx.NewNode("Deformer", id, a2, a3) }
func Document(id int64, a2, a3 string) *fbx.Node      { return fbx.NewNode("Document", id, a2, a3) }
func Documents() *fbx.Node                            { return fbx.NewNode("Documents") }
func EncryptionType(t int32) *fbx.Node                { return fbx.NewNode("EncryptionType", t) }
func FBXHeaderExtension() *fbx.Node                   { return fbx.NewNode("FBXHeaderExtension") }
func FBXHeaderVersion(version int32) *fbx.Node        { return fbx.NewNode("FBXHeaderVersion", version) }
func FBXVersion(version int32) *fbx.Node              { return fbx.NewNode("FBXVersion", version) }
func FileId(id []byte) *fbx.Node                      { return fbx.NewNode("FileId", id) }
func Filename(filename string) *fbx.Node              { return fbx.NewNode("Filename", filename) }
func Geometry(id int64, a2, a3 string) *fbx.Node      { return fbx.NewNode("Geometry", id, a2, a3) }
func GeometryVersion(version int32) *fbx.Node         { return fbx.NewNode("GeometryVersion", version) }
func GlobalSettings() *fbx.Node                       { return fbx.NewNode("GlobalSettings") }
func Hour(hour int32) *fbx.Node                       { return fbx.NewNode("Hour", hour) }
func Keywords(keywords string) *fbx.Node              { return fbx.NewNode("Keywords", keywords) }
func Layer(a1 int32) *fbx.Node                        { return fbx.NewNode("Layer", a1) }
func LayerElement() *fbx.Node                         { return fbx.NewNode("LayerElement") }
func LayerElementMaterial(a1 int32) *fbx.Node         { return fbx.NewNode("LayerElementMaterial", a1) }
func LayerElementNormal(a1 int32) *fbx.Node           { return fbx.NewNode("LayerElementNormal", a1) }
func LayerElementColor(a1 int32) *fbx.Node            { return fbx.NewNode("LayerElementColor", a1) }
func LayerElementUV(a1 int32) *fbx.Node               { return fbx.NewNode("LayerElementUV", a1) }
func Link_DeformAcuracy(a1 float64) *fbx.Node         { return fbx.NewNode("Link_DeformAcuracy", a1) }
func MappingInformationType(t string) *fbx.Node       { return fbx.NewNode("MappingInformationType", t) }
func Material(id int64, a2, a3 string) *fbx.Node      { return fbx.NewNode("Material", id, a2, a3) }
func Materials(a1 []int32) *fbx.Node                  { return fbx.NewNode("Materials", a1) }
func Matrix(matrix []float64) *fbx.Node               { return fbx.NewNode("Matrix", matrix) }
func Media(a1 string) *fbx.Node                       { return fbx.NewNode("Media", a1) }
func MetaData() *fbx.Node                             { return fbx.NewNode("MetaData") }
func Millisecond(millisecond int32) *fbx.Node         { return fbx.NewNode("Millisecond", millisecond) }
func Minute(minute int32) *fbx.Node                   { return fbx.NewNode("Minute", minute) }
func Model(id int64, a2, a3 string) *fbx.Node         { return fbx.NewNode("Model", id, a2, a3) }
func ModelUVScaling(a1, a2 float64) *fbx.Node         { return fbx.NewNode("ModelUVScaling", a1, a2) }
func ModelUVTranslation(a1, a2 float64) *fbx.Node     { return fbx.NewNode("ModelUVTranslation", a1, a2) }
func Month(month int32) *fbx.Node                     { return fbx.NewNode("Month", month) }
func MultiLayer(a1 int32) *fbx.Node                   { return fbx.NewNode("MultiLayer", a1) }
func MultiTake(a1 int32) *fbx.Node                    { return fbx.NewNode("MultiTake", a1) }
func Name(name string) *fbx.Node                      { return fbx.NewNode("Name", name) }
func NbPoseNodes(count int32) *fbx.Node               { return fbx.NewNode("NbPoseNodes", count) }
func Node(id int64) *fbx.Node                         { return fbx.NewNode("Node", id) }
func NodeAttribute(id int64, a2, a3 string) *fbx.Node { return fbx.NewNode("NodeAttribute", id, a2, a3) }
func Normals(normals []float64) *fbx.Node             { return fbx.NewNode("Normals", normals) }
func Indexes(indexes []int32) *fbx.Node               { return fbx.NewNode("Indexes", indexes) }
func Objects() *fbx.Node                              { return fbx.NewNode("Objects") }
func ObjectType(object string) *fbx.Node              { return fbx.NewNode("ObjectType", object) }
func PolygonVertexIndex(indexes []int32) *fbx.Node    { return fbx.NewNode("PolygonVertexIndex", indexes) }
func Pose(id int64, a2, a3 string) *fbx.Node          { return fbx.NewNode("Pose", id, a2, a3) }
func PoseNode() *fbx.Node                             { return fbx.NewNode("PoseNode") }
func Properties70() *fbx.Node                         { return fbx.NewNode("Properties70") }
func PropertyTemplate(a1 string) *fbx.Node            { return fbx.NewNode("PropertyTemplate", a1) }
func ReferenceInformationType(t string) *fbx.Node     { return fbx.NewNode("ReferenceInformationType", t) }
func References() *fbx.Node                           { return fbx.NewNode("References") }
func RelativeFilename(filename string) *fbx.Node      { return fbx.NewNode("RelativeFilename", filename) }
func Revision(revision string) *fbx.Node              { return fbx.NewNode("Revision", revision) }
func RootNode(id int64) *fbx.Node                     { return fbx.NewNode("RootNode", id) }
func SceneInfo(a1 string, a2 string) *fbx.Node        { return fbx.NewNode("SceneInfo", a1, a2) }
func Second(second int32) *fbx.Node                   { return fbx.NewNode("Second", second) }
func Shading(shading bool) *fbx.Node                  { return fbx.NewNode("Shading", shading) }
func ShadingModel(model string) *fbx.Node             { return fbx.NewNode("ShadingModel", model) }
func SkinningType(t string) *fbx.Node                 { return fbx.NewNode("SkinningType", t) }
func Subject(subject string) *fbx.Node                { return fbx.NewNode("Subject", subject) }
func Takes() *fbx.Node                                { return fbx.NewNode("Takes") }
func Texture(id int64, a2, a3 string) *fbx.Node       { return fbx.NewNode("Texture", id, a2, a3) }
func Texture_Alpha_Source(a1 string) *fbx.Node        { return fbx.NewNode("Texture_Alpha_Source", a1) }
func TextureName(name string) *fbx.Node               { return fbx.NewNode("TextureName", name) }
func Title(title string) *fbx.Node                    { return fbx.NewNode("Title", title) }
func Transform(matrix []float64) *fbx.Node            { return fbx.NewNode("Transform", matrix) }
func TransformLink(matrix []float64) *fbx.Node        { return fbx.NewNode("TransformLink", matrix) }
func Type(a1 string) *fbx.Node                        { return fbx.NewNode("Type", a1) }
func TypedIndex(index int32) *fbx.Node                { return fbx.NewNode("TypedIndex", index) }
func TypeFlags(t string) *fbx.Node                    { return fbx.NewNode("TypeFlags", t) }
func UseMipMap(a1 int32) *fbx.Node                    { return fbx.NewNode("UseMipMap", a1) }
func UserData(a1, a2 string) *fbx.Node                { return fbx.NewNode("UserData", a1, a2) }
func UV(vertices []float64) *fbx.Node                 { return fbx.NewNode("UV", vertices) }
func UVIndex(indexes []int32) *fbx.Node               { return fbx.NewNode("UVIndex", indexes) }
func Version(version int32) *fbx.Node                 { return fbx.NewNode("Version", version) }
func Vertices(vertices []float64) *fbx.Node           { return fbx.NewNode("Vertices", vertices) }
func Video(id int64, a2, a3 string) *fbx.Node         { return fbx.NewNode("Video", id, a2, a3) }
func Weights(weights []float64) *fbx.Node             { return fbx.NewNode("Weights", weights) }
func Year(year int32) *fbx.Node                       { return fbx.NewNode("Year", year) }
