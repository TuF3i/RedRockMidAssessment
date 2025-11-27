package snowflake

import (
	"RedRockMidAssessment/core"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"sync"
	"time"
)

//  0                   1                   2                   3
//  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |0|000 0000 0000 0000 0000 0000 0000 0000 0000 0000 00|00 0000 0000|0000 0000 0000|
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |0|  41-bit 时间差                                     | 10-bit 机器 | 12-bit 序号   |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// 第一位符号位，始终为0

const (
	epoch        = int64(1764000000000) // 发号器元年 2025-11-25 00:00:00 UTC
	machineBits  = 10                   // 机器位位数，支持1024台机器
	sequenceBits = 12                   // 每毫秒单台机器可发 4096 个号

	machineMax   = -1 ^ (-1 << machineBits)  // 1111111111111111111111111111111111111111111111111111111111111111 xor 1111111111111111111111111111111111111111111111111111110000000000
	sequenceMask = -1 ^ (-1 << sequenceBits) // 1111111111111111111111111111111111111111111111111111111111111111 xor 1111111111111111111111111111111111111111111111111111000000000000
	maxTimestamp = -1 ^ (-1 << (41))         // 69 年
)

func NewSnowflake(machineID int64) error {
	// 初始化Snowflake，machineID可用 IP 末段，K8S Pod 序号
	if machineID < 0 || machineID > machineMax {
		return errors.New("machine id out of range")
	}

	core.SnowFlake = &Snowflake{machineID: machineID}

	return nil
}

type Snowflake struct {
	sync.Mutex       // 锁，防止同一毫秒下两个G竞争一个Snowflake
	machineID  int64 // 机器ID
	sequence   int64 // 当前毫秒已发多少号（0~4095）
	lastTime   int64 // 上一次发号的时间戳（毫秒）
}

func (s *Snowflake) gen() int64 {
	s.Lock()
	defer s.Unlock()

	now := time.Now().UnixMilli()

	// ---------- 时钟回拨处理 ----------
	if now < s.lastTime {
		drift := s.lastTime - now
		for now < s.lastTime { // 循环直到追上
			time.Sleep(time.Duration(drift+1) * time.Millisecond)
			now = time.Now().UnixMilli()
			drift = s.lastTime - now
			if drift > 1000 { // 回拨超过 1s 直接报错，避免无限睡
				// panic(fmt.Sprintf("clock moved backwards %d ms", drift))
				return s.genUUID() // 降级用UUID
			}
		}
	}
	// ---------- 后续逻辑不变 ----------
	if now == s.lastTime {
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			for now <= s.lastTime {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		s.sequence = 0
	}
	s.lastTime = now
	id := (now-epoch)<<22 | s.machineID<<12 | s.sequence
	return id
}

func (s *Snowflake) TraceID() string {
	return fmt.Sprintf("%016x", s.gen())
}

func (s *Snowflake) genUUID() int64 {
	u := uuid.New()
	bytes := u[:8]
	return int64(binary.BigEndian.Uint64(bytes))
}
