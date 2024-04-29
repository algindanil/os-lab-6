package main

import (
	"C"
	"github.com/remeh/sizedwaitgroup"
	"gonum.org/v1/gonum/mat"
	"math/rand"
	"os"
	"time"
)

func GetRandMatrix(rows, cols int) *mat.Dense {
	data := make([]float64, rows*cols)
	for i := range data {
		data[i] = rand.NormFloat64()
	}
	return mat.NewDense(rows, cols, data)
}

func CalcElement(resMatrix *mat.Dense, wg *sizedwaitgroup.SizedWaitGroup, leftVec, rightVec mat.Vector, m, n int) {
	newElement := mat.Dot(leftVec, rightVec)
	resMatrix.Set(m, n, newElement)

	//fmt.Printf("Element at (%d;%d): %.3f\n", m, n, newElement)

	wg.Done()
}

func ParallelMatMul(left, right *mat.Dense, maxThreads int) *mat.Dense {
	m, _ := left.Dims()
	_, n := right.Dims()

	resMatrix := mat.NewDense(m, n, nil)

	wg := sizedwaitgroup.New(maxThreads)

	for i := range m {
		for j := range n {
			leftVec := left.RowView(i)
			rightVec := right.ColView(j)
			wg.Add()
			go CalcElement(resMatrix, &wg, leftVec, rightVec, i, j)
		}
	}

	wg.Wait()
	return resMatrix
}

//export TimeMatMul
func TimeMatMul(m, k, n, maxThreads int) int {
	left := GetRandMatrix(m, k)
	right := GetRandMatrix(k, n)

	start := time.Now()
	_ = ParallelMatMul(left, right, maxThreads)
	elapsed := time.Since(start).Milliseconds()

	result := int(elapsed)
	return result
}

func main() {
	res := TimeMatMul(1000, 1000, 1000, 1000)
	os.Exit(res)
}
