package offsetboundary

import (
	"strings"
	"testing"

	"github.com/relunctance/goutils/slice"
	"github.com/smartystreets/goconvey/convey"
)

var data [][]string = [][]string{
	[]string{"1a", "1b", "1c"},                                                 //3
	[]string{"2a", "2b", "2c", "2d", "2e", "2f", "2g"},                         //7
	[]string{"3a", "3b", "3c", "3d", "3e", "3f", "3g", "3h", "3i"},             //9
	[]string{"4a", "4b", "4c", "4d", "4e", "4f", "4g", "4h", "4i", "4j", "4k"}, //11
}

//获取数据结果
func getData(arr [][]int) (res []string) {

	for key, items := range arr {
		if items == nil { //这里如果为nil 说明是跳过的
			continue
		}
		start, end := GetStartEndByItems(items)
		res = append(res, data[key][start:end]...)
		/*
			for _, v := range items {	//效率不好
				res = append(res, data[key][v])
			}
		*/
	}
	return
}
func TestMultiIndexDataHasNext(t *testing.T) {
	convey.Convey("测试二维索引数组是否结束", t, func() {
		convey.Convey("测试返回bool:", func() {
			indexData := [][]int{
				{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14},
				{0, 1, 2, 3, 4, 5, 6, 7},
			}
			var offsetIndexData [][]int
			offsetIndexData = [][]int{{1, 2, 3}, {5}}
			convey.So(MultiIndexDataHasNext(indexData, offsetIndexData), convey.ShouldBeTrue)
			offsetIndexData = [][]int{{}, {6}}
			convey.So(MultiIndexDataHasNext(indexData, offsetIndexData), convey.ShouldBeTrue)
			offsetIndexData = [][]int{{}, {5, 6, 7}}
			convey.So(MultiIndexDataHasNext(indexData, offsetIndexData), convey.ShouldBeFalse)
			offsetIndexData = [][]int{{14}, {}}
			convey.So(MultiIndexDataHasNext(indexData, offsetIndexData), convey.ShouldBeFalse)
			offsetIndexData = [][]int{{}, {}}
			convey.So(MultiIndexDataHasNext(indexData, offsetIndexData), convey.ShouldBeFalse)
		})

		convey.Convey("测试初始值中indexData存在空的情况:", func() {

			indexData := [][]int{
				{},
				{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14},
				{},
				{0, 1, 2, 3, 4, 5, 6, 7},
			}
			var offsetIndexData [][]int
			offsetIndexData = [][]int{{}, {1, 2, 3}, {}, {3, 4}}
			convey.So(MultiIndexDataHasNext(indexData, offsetIndexData), convey.ShouldBeTrue)
			offsetIndexData = [][]int{{}, {1, 2, 3}, {}, {7}} //中间的{1,2,3} 已经直接导致返回了true
			convey.So(MultiIndexDataHasNext(indexData, offsetIndexData), convey.ShouldBeTrue)

			offsetIndexData = [][]int{{}, {14}, {}, {3, 4}} //{3,4} 导致返回了true
			convey.So(MultiIndexDataHasNext(indexData, offsetIndexData), convey.ShouldBeTrue)

			offsetIndexData = [][]int{{}, {}, {}, {7}}
			convey.So(MultiIndexDataHasNext(indexData, offsetIndexData), convey.ShouldBeFalse)

			offsetIndexData = [][]int{{}, {14}, {}, {7}}
			convey.So(MultiIndexDataHasNext(indexData, offsetIndexData), convey.ShouldBeFalse)

			offsetIndexData = [][]int{{}, {}, {1, 2, 3}, {}} //只要有一个len == 0 的情况下直接跳过的
			convey.So(MultiIndexDataHasNext(indexData, offsetIndexData), convey.ShouldBeFalse)

		})

		convey.Convey("测试panic报错:", func() {

			indexData := [][]int{
				{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14},
				{0, 1, 2, 3, 4, 5, 6, 7},
			}
			offsetIndexData := [][]int{
				{1, 2, 3},
				{5},
				{1, 2, 3},
			}
			convey.So(func() {
				MultiIndexDataHasNext(indexData, offsetIndexData)
			}, convey.ShouldPanic)
		})
	})

}

func TestIsEnd(t *testing.T) {
	convey.Convey("测试是否offsetdata是否结束:", t, func() {
		arr := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
		isend, offset := CheckIsOffsetEnd(arr, []int{12, 13, 14})
		convey.So(isend, convey.ShouldBeTrue)
		convey.So(offset, convey.ShouldEqual, 0)
		isend, offset = CheckIsOffsetEnd(arr, []int{7, 8, 9, 10})
		convey.So(isend, convey.ShouldBeFalse)
		convey.So(offset, convey.ShouldEqual, 11)

		isend, offset = CheckIsOffsetEnd(arr, []int{14})
		convey.So(isend, convey.ShouldBeTrue)
		convey.So(offset, convey.ShouldEqual, 0)

		isend, offset = CheckIsOffsetEnd(arr, []int{0, 1, 2, 3, 4, 5})
		convey.So(isend, convey.ShouldBeFalse)
		convey.So(offset, convey.ShouldEqual, 6)
		isend, offset = CheckIsOffsetEnd(arr, arr)
		convey.So(isend, convey.ShouldBeTrue)
		convey.So(offset, convey.ShouldEqual, 0)
		//检测panic
		convey.So(func() { CheckIsOffsetEnd(arr, append(arr, 15)) }, convey.ShouldPanic)

		isend, offset = CheckIsOffsetEnd([]int{1, 2, 3, 4, 5, 6, 7}, []int{3, 4})
		convey.So(isend, convey.ShouldBeFalse)
		convey.So(offset, convey.ShouldEqual, 5)

		isend, offset = CheckIsOffsetEnd([]int{1, 2, 3, 4, 5, 6, 7}, []int{})
		convey.So(isend, convey.ShouldBeFalse)
		convey.So(offset, convey.ShouldEqual, 0)

		isend, offset = CheckIsOffsetEnd([]int{}, []int{})
		convey.So(isend, convey.ShouldBeFalse)
		convey.So(offset, convey.ShouldEqual, 0)
	})
}

func TestGetDyadicArrayByOffset2(t *testing.T) {
	convey.Convey("测试单独二维slice", t, func() {
		multiDyadicArr := [][]int{
			{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14},
			{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		}
		totalLength := int64(len(multiDyadicArr[0]) + len(multiDyadicArr[1]))
		offset := int64(8)
		pagesize := int64(10)
		_, arr := GetDyadicArrayByOffset(offset, pagesize, totalLength, multiDyadicArr)
		convey.Convey("数组1应该相等:", func() {
			b1 := slice.CheckIntSliceEqual(arr[0], []int{8, 9, 10, 11, 12, 13, 14})
			convey.So(b1, convey.ShouldBeTrue)
			b2 := slice.CheckIntSliceEqual(arr[1], []int{0, 1, 2})
			convey.So(b2, convey.ShouldBeTrue)
		})
	})
}

func TestGetDyadicArrayByOffset(t *testing.T) {
	convey.Convey("测试二维数组的offset,pagesize检索", t, func() {
		multiDyadicArr, totalLength := BuildDyadicArray(data) //totalLength == 30

		offset := int64(0)
		pagesize := int64(10)
		err, arr := GetDyadicArrayByOffset(offset, pagesize, totalLength, multiDyadicArr)
		convey.Convey("err应该返回nil", func() {
			convey.So(nil, convey.ShouldEqual, err)
		})
		convey.Convey("arr值判断1:", func() {
			bl := slice.CheckStringSliceEqual(strings.Fields("1a 1b 1c 2a 2b 2c 2d 2e 2f 2g"), getData(arr))
			convey.So(bl, convey.ShouldBeTrue)
		})

		err, arr = GetDyadicArrayByOffset(10000, pagesize, totalLength, multiDyadicArr)
		convey.Convey("检测err,应该不为空", func() {
			convey.So(err, convey.ShouldNotBeNil)
		})
		err, arr = GetDyadicArrayByOffset(30, pagesize, totalLength, multiDyadicArr)
		convey.Convey("检测结束符号", func() {
			convey.So(err, convey.ShouldEqual, EOF)
		})

		err, arr = GetDyadicArrayByOffset(5, 10, totalLength, multiDyadicArr)
		convey.Convey("arr值判断2:", func() {
			bl := slice.CheckStringSliceEqual(strings.Fields("2c 2d 2e 2f 2g 3a 3b 3c 3d 3e"), getData(arr))
			convey.So(bl, convey.ShouldBeTrue)
			convey.So(err, convey.ShouldBeNil)
		})
		err, arr = GetDyadicArrayByOffset(3, 5, totalLength, multiDyadicArr)
		convey.Convey("检测当start为边界时 , 是否报错问题:", func() {
			bl := slice.CheckStringSliceEqual(strings.Fields("2a 2b 2c 2d 2e"), getData(arr))
			convey.So(bl, convey.ShouldBeTrue)
			convey.So(err, convey.ShouldBeNil)

		})
		err, arr = GetDyadicArrayByOffset(0, 5000000, totalLength, multiDyadicArr)
		convey.Convey("当pagesize过大时检测", func() {
			bl := slice.CheckStringSliceEqual(strings.Fields("1a 1b 1c 2a 2b 2c 2d 2e 2f 2g 3a 3b 3c 3d 3e 3f 3g 3h 3i 4a 4b 4c 4d 4e 4f 4g 4h 4i 4j 4k"), getData(arr))
			convey.So(bl, convey.ShouldBeTrue)
			convey.So(err, convey.ShouldBeNil)

		})

		err, arr = GetDyadicArrayByOffset(13, 5000000, totalLength, multiDyadicArr)
		convey.Convey("检测", func() {
			bl := slice.CheckStringSliceEqual(strings.Fields("3d 3e 3f 3g 3h 3i 4a 4b 4c 4d 4e 4f 4g 4h 4i 4j 4k"), getData(arr))
			convey.So(bl, convey.ShouldBeTrue)
			convey.So(err, convey.ShouldBeNil)
			convey.So(len(getData(arr)), convey.ShouldEqual, int(totalLength-13))
		})
	})
}

func TestBuildDyadicArray(t *testing.T) {
	convey.Convey("测试多维数组的搜索", t, func() {
		arr := [][]string{}
		mutliarr, num := BuildDyadicArray(arr)
		convey.Convey("长度应该为0", func() {
			convey.So(0, convey.ShouldEqual, len(mutliarr))
		})
		convey.Convey("结果应该为0", func() {
			convey.So(num, convey.ShouldEqual, 0)
		})

		mutliarr, num = BuildDyadicArray(data)

		convey.Convey("结果应该相等:数组长度", func() {
			convey.So(len(mutliarr), convey.ShouldEqual, len(data))
		})
		convey.Convey("结果应该相等:总长", func() {
			convey.So(30, convey.ShouldEqual, num)
		})
	})

}
