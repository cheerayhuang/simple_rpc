package main

import (
	"bytes"
	"encoding/gob"
	"reflect"
	"math"
	//"fmt"

	"github.com/rs/zerolog"
)

type Vector struct {
	X, Y, Z int
}

func testGob() {
	var vBuffer bytes.Buffer	

	// ex. en/decode a struct
	enc := gob.NewEncoder(&vBuffer)
	if err := enc.Encode(Vector{3,4,5}); err != nil {
		logger.Fatal().Err(err).Msg("encode Vector to gob failed")
	}

	var v = new(Vector)


	dec := gob.NewDecoder(&vBuffer)
	if err := dec.Decode(&v); err != nil {
		logger.Fatal().Err(err).Msg("decode gob to Vector failed")
	}

	logger.Info().Msgf("Vector after decode: %v", v)

	// ex. en/decode a num
	enc.Encode(5)

	var floatNum float32
	dec.Decode(&floatNum)

	logger.Info().Msgf("int after decode into a float: %v", floatNum) // wrong, output: 0

	enc.Encode(5.0)
	dec.Decode(&floatNum)
	logger.Info().Msgf("float after decode into a float: %v", floatNum) // right, output: 5 

	var floatNum2 float64
	enc.Encode(5.0)
	dec.Decode(&floatNum2)
	logger.Info().Msgf("float after decode into a float: %v", floatNum2) // right, output: 5 

	var intNum int32
	enc.Encode(5)
	dec.Decode(&intNum)
	logger.Info().Msgf("int after decode into a int: %v", intNum) // right, output: 5 
}

func testZerolog() {
	var v = Vector{10, 9, 8}
	var pointerV = &Vector{3, 4, 5}

	logger.Info().Any("vector_any", v).Object("vector_obj", v).EmbedObject(v).Send()
	logger.Info().EmbedObject(v).Send()

	v.Handler1(11)
	logger.Info().EmbedObject(v).Send()
	
	v.Handler2(12)
	logger.Info().EmbedObject(v).Send()

	pointerV.Handler1(30)
	logger.Info().EmbedObject(pointerV).Send()

	pointerV.Handler2(40)
	logger.Info().EmbedObject(pointerV).Send()
}

func testReflect() {

	var v = Vector{3, 4, 5}
	var pV = new(Vector)
	
	typeV := reflect.TypeOf(v)
	typePV := reflect.TypeOf(pV)
	valV := reflect.ValueOf(v)

	logger.Info().Str("type", typeV.Name()).Int("method_num", typeV.NumMethod()).Send()
	logger.Info().Str("type of pointer to V", typePV.String()).Send() // pointer can't be displayed by Name()

	for i := 0; i < typeV.NumMethod(); i++ {
		method := typeV.Method(i)
		typeM := method.Type

		logger.Info().Str(method.Name, typeM.String()).Send()

		if method.Name == "Handler2" {
			method.Func.Call([]reflect.Value{
				valV,
				reflect.ValueOf(float32(5.0)),	
			})
		
		}
	}
}

func (v *Vector) Handler1(val int) {
	v.X = val 
}

func (v Vector) Handler2(val float32) {
	v.Y = int(val)

	logger.Info().Float32("try_to_set", val).Msg("Running in Handler2.")
}

func (v Vector) MarshalZerologObject(e *zerolog.Event) {
	e.Int("x", v.X).Int("Y", v.Y).Int("Z", v.Z).Float64("D", math.Sqrt(float64(v.X*v.X + v.Y*v.Y+v.Z*v.Z)))
}

func complexFunc(a int32, b float32, c string) {
	logger.Info().
		Int32("a", a).Float32("b", b).Str("c", c).Msg("Running in complexFunc")
}

type Req struct {
	paramsTypes []reflect.Type
	paramsInitVal []byte
	vectorsInitVal []byte
	paramsEachVal [][]byte
	funcVal reflect.Value
}

func testUseGobToCallFunc() {

	argVec := Vector{3, 4, 5}	
	a := int32(1000)	
	b := float32(3.14159)
	c := string("hello world")
	
	//args := []any{&argVec, a, b, c}
	args := []any{a, b, c}
	argVecs := []Vector{argVec}

	argsBuf := new(bytes.Buffer)
	enc := gob.NewEncoder(argsBuf)
	dec := gob.NewDecoder(argsBuf)
	enc.Encode(args)

	var args2 []any 
	dec.Decode(&args2)
	logger.Info().Any("args", args).Send()
	logger.Info().Any("decode_args", args2).Send()

	var argVecs2 []Vector
	enc.Encode(argVecs)
	dec.Decode(&argVecs2)
	logger.Info().Any("argVecs", argVecs).Send()
	logger.Info().Any("decode_argVecs", argVecs2).Send()

	enc.Encode(args)
	argVecsBuf := new(bytes.Buffer)
	encVecs := gob.NewEncoder(argVecsBuf) 
	encVecs.Encode(argVecs)

	req := &Req {
		paramsInitVal: argsBuf.Bytes(),
		vectorsInitVal: argVecsBuf.Bytes(),
	}
	logger.Info().Int("buf_len_req", len(req.vectorsInitVal)).Send()

	// any 类型的数组没有办法在通过函数调用之后再decode回来
	callDecode(argsBuf.Bytes())

	// 但是[]Vector类型，是一个明确的slice类型，因此可以
	RPCSimulateInvokeViaVals(req)

	req.funcVal = reflect.ValueOf(complexFunc)
	funcType := reflect.TypeOf(complexFunc)
	logger.Info().Int("num_IN", funcType.NumIn()).Send()
	for i := 0; i < funcType.NumIn(); i++ {
		req.paramsTypes = append(req.paramsTypes, funcType.In(i))
	}
	
	tmpBuf := new(bytes.Buffer)
	encEach := gob.NewEncoder(tmpBuf)
	encEach.Encode(a)
	req.paramsEachVal = append(req.paramsEachVal, tmpBuf.Bytes())

	appendBytes()

	// 这里如果使用同样的buffer 和encoder，则会造成连续encode
	// 产生 appendBytes() 里描述的行为
	tmpBufB := new(bytes.Buffer)
	encB := gob.NewEncoder(tmpBufB)
	encB.Encode(b)
	req.paramsEachVal = append(req.paramsEachVal, tmpBufB.Bytes())

	tmpBufC := new(bytes.Buffer)
	encC := gob.NewEncoder(tmpBufC)
	encC.Encode(c)
	req.paramsEachVal = append(req.paramsEachVal, tmpBufC.Bytes())

	logger.Info().Any("each_val", req.paramsEachVal).Send()

	RPCSimulateInvokeViaTypes(req)
}

func appendBytes() {
	tmpBytes := []byte{1}
	tmpBuf := bytes.NewBuffer(tmpBytes)

	var a [][]byte

	a = append(a, tmpBuf.Bytes())
	a = append(a, tmpBuf.Bytes())
	logger.Info().Any("a", a).Send()

	enc := gob.NewEncoder(tmpBuf)
	//tmpBytes[0] = 2
	enc.Encode(2)
	logger.Info().Any("a", a).Send()
	a = append(a, tmpBuf.Bytes())
	logger.Info().Any("a", a).Send()

	var tmpBuf2 bytes.Buffer
	enc2 := gob.NewEncoder(&tmpBuf2)
	enc2.Encode(3)
	a = append(a, tmpBuf2.Bytes())
	logger.Info().Any("a", a).Send()
	// 连续 encode 的行为很神奇，它会重新new一块bytes并且复制之前的bytes的值。
	enc2.Encode(4)
	a = append(a, tmpBuf2.Bytes())
	logger.Info().Any("a", a).Send()
}

func callDecode(buf []byte) {
	argsBuf := bytes.NewBuffer(buf)
	dec := gob.NewDecoder(argsBuf)

	var args *[]any

	err := dec.Decode(&args)
	if err != nil {
		logger.Error().Err(err).Any("args", args).Msg("decode []any in another func failed.")
		return
	}

	logger.Info().Any("args", args).Send()
}

func RPCSimulateInvokeViaVals(req *Req) {
	// recover to get args' value
	logger.Info().Int("buf_len", len(req.vectorsInitVal)).Send()
	vecsBuf := bytes.NewBuffer(req.vectorsInitVal)
	dec := gob.NewDecoder(vecsBuf)
	var vecs []Vector
	if err := dec.Decode(&vecs); err != nil {
		logger.Error().Err(err).Msg("decode []Vector in another func failed.")
		return
	}

	logger.Info().Any("vecs", vecs).Send()
	//fmt.Println(args)

	/*
	funcVal := reflect.ValueOf(complexFunc)
	params := make([]reflect.Value, 4)
	for i, v := range argsVals {
		params[i] = reflect.ValueOf(v)
	}
	funcVal.Call(params)
	*/
}


func RPCSimulateInvokeViaTypes(req *Req) {
	var args []reflect.Value

	for i, b := range req.paramsEachVal {
		arg := reflect.New(req.paramsTypes[i])
		buf := bytes.NewBuffer(b)
		dec := gob.NewDecoder(buf)
		dec.Decode(arg.Interface())

		logger.Info().Str("kind", arg.Elem().String()).Send()
		args = append(args, arg.Elem())
	}

	logger.Info().Any("args", args).Send()

	req.funcVal.Call(args)
}

func main() {
	logger.Info().Msg("this is a test program.")
	//logger.Info().Int("a", 4).Send()
	
	testGob()

	testZerolog()

	testReflect()

	testUseGobToCallFunc()
}
