/*
 * @Author: webees@qq.com
 * @Date: 2021-03-29 18:10:10
 * @Last Modified by:   webees@qq.com
 * @Last Modified time: 2021-03-29 18:10:10
 */
package fil

func coinType() uint32 {
	if TEST {
		return tCoinID
	} else {
		return coinID
	}
}
