/**
* @Author: Cooper
* @Date: 2019/12/4 20:21
 */

package mface

type MConnection interface {
	MBase

	ID() string
	ResponseChan() chan MMessage
	Status() int
}
