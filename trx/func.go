/*
 * @Author: webees@qq.com
 * @Date: 2021-03-29 18:09:55
 * @Last Modified by: webees@qq.com
 * @Last Modified time: 2021-03-29 19:04:45
 */
package trx

func coinType() uint32 {
	if TEST {
		return tCoinID
	} else {
		return coinID
	}
}
