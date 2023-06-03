package metrics

type Unit string

// Copied from https://docs.datadoghq.com/metrics/units/#unit-list
const (
	// BYTES
	UnitBit      Unit = "bit"
	UnitByte     Unit = "byte"
	UnitKibibyte Unit = "kibibyte"
	UnitMebibyte Unit = "mebibyte"
	UnitGibibyte Unit = "gibibyte"
	UnitTebibyte Unit = "tebibyte"
	UnitPebibyte Unit = "pebibyte"
	UnitExbibyte Unit = "exbibyte"

	// TIME
	UnitNanosecond  Unit = "nanosecond"
	UnitMicrosecond Unit = "microsecond"
	UnitMillisecond Unit = "millisecond"
	UnitSecond      Unit = "second"
	UnitMinute      Unit = "minute"
	UnitHour        Unit = "hour"
	UnitDay         Unit = "day"
	UnitWeek        Unit = "week"

	// PERCENTAGE
	UnitPercent_nano Unit = "percent_nano"
	UnitPercent      Unit = "percent"
	UnitApdex        Unit = "apdex"
	UnitFraction     Unit = "fraction"

	// NETWORK
	UnitConnection Unit = "connection"
	UnitRequest    Unit = "request"
	UnitPacket     Unit = "packet"
	UnitSegment    Unit = "segment"
	UnitResponse   Unit = "response"
	UnitMessage    Unit = "message"
	UnitPayload    Unit = "payload"
	UnitTimeout    Unit = "timeout"
	UnitDatagram   Unit = "datagram"
	UnitRoute      Unit = "route"
	UnitSession    Unit = "session"
	UnitHop        Unit = "hop"

	// SYSTEM
	UnitProcess  Unit = "process"
	UnitThread   Unit = "thread"
	UnitHost     Unit = "host"
	UnitNode     Unit = "node"
	UnitFault    Unit = "fault"
	UnitService  Unit = "service"
	UnitInstance Unit = "instance"
	UnitCpu      Unit = "cpu"

	// DISK
	UnitFile   Unit = "file"
	UnitInode  Unit = "inode"
	UnitSector Unit = "sector"
	UnitBlock  Unit = "block"

	// GENERAL
	UnitBuffer            Unit = "buffer"
	UnitError             Unit = "error"
	UnitRead              Unit = "read"
	UnitWrite             Unit = "write"
	UnitOccurrence        Unit = "occurrence"
	UnitEvent             Unit = "event"
	UnitTime              Unit = "time"
	UnitUnit              Unit = "unit"
	UnitOperation         Unit = "operation"
	UnitItem              Unit = "item"
	UnitTask              Unit = "task"
	UnitWorker            Unit = "worker"
	UnitResource          Unit = "resource"
	UnitGarbageCollection Unit = "garbage collection"
	UnitEmail             Unit = "email"
	UnitSample            Unit = "sample"
	UnitStage             Unit = "stage"
	UnitMonitor           Unit = "monitor"
	UnitLocation          Unit = "location"
	UnitCheck             Unit = "check"
	UnitAttempt           Unit = "attempt"
	UnitDevice            Unit = "device"
	UnitUpdate            Unit = "update"
	UnitMethod            Unit = "method"
	UnitJob               Unit = "job"
	UnitContainer         Unit = "container"
	UnitExecution         Unit = "execution"
	UnitThrottle          Unit = "throttle"
	UnitInvocation        Unit = "invocation"
	UnitUser              Unit = "user"
	UnitSuccess           Unit = "success"
	UnitBuild             Unit = "build"
	UnitPrediction        Unit = "prediction"
	UnitException         Unit = "exception"

	// DB
	UnitTable       Unit = "table"
	UnitIndex       Unit = "index"
	UnitLock        Unit = "lock"
	UnitTransaction Unit = "transaction"
	UnitQuery       Unit = "query"
	UnitRow         Unit = "row"
	UnitKey         Unit = "key"
	UnitCommand     Unit = "command"
	UnitOffset      Unit = "offset"
	UnitRecord      Unit = "record"
	UnitObject      Unit = "object"
	UnitCursor      Unit = "cursor"
	UnitAssertion   Unit = "assertion"
	UnitScan        Unit = "scan"
	UnitDocument    Unit = "document"
	UnitShard       Unit = "shard"
	UnitFlush       Unit = "flush"
	UnitMerge       Unit = "merge"
	UnitRefresh     Unit = "refresh"
	UnitFetch       Unit = "fetch"
	UnitColumn      Unit = "column"
	UnitCommit      Unit = "commit"
	UnitWait        Unit = "wait"
	UnitTicket      Unit = "ticket"
	UnitQuestion    Unit = "question"

	// CACHE
	UnitHit      Unit = "hit"
	UnitMiss     Unit = "miss"
	UnitEviction Unit = "eviction"
	UnitGet      Unit = "get"
	UnitSet      Unit = "set"

	// MONEY
	UnitDollar      Unit = "dollar"
	UnitCent        Unit = "cent"
	UnitMicrodollar Unit = "microdollar"
	UnitEuro        Unit = "euro"

	// MEMORY
	UnitPage  Unit = "page"
	UnitSplit Unit = "split"

	// FREQUENCY
	UnitHertz     Unit = "hertz"
	UnitKilohertz Unit = "kilohertz"
	UnitMegahertz Unit = "megahertz"
	UnitGigahertz Unit = "gigahertz"

	// LOGGING
	UnitEntry Unit = "entry"

	// TEMPERATURE
	UnitDecidegreeCelsius Unit = "decidegree celsius"
	UnitDegreeCelsius     Unit = "degree celsius"
	UnitDegreeFahrenheit  Unit = "degree fahrenheit"

	// CPU
	UnitNanocore  Unit = "nanocore"
	UnitMicrocore Unit = "microcore"
	UnitMillicore Unit = "millicore"
	UnitCore      Unit = "core"
	UnitKilocore  Unit = "kilocore"
	UnitMegacore  Unit = "megacore"
	UnitGigacore  Unit = "gigacore"
	UnitTeracore  Unit = "teracore"
	UnitPetacore  Unit = "petacore"
	UnitExacore   Unit = "exacore"

	// POWER
	UnitNanowatt  Unit = "nanowatt"
	UnitMicrowatt Unit = "microwatt"
	UnitMilliwatt Unit = "milliwatt"
	UnitDeciwatt  Unit = "deciwatt"
	UnitWatt      Unit = "watt"
	UnitKilowatt  Unit = "kilowatt"
	UnitMegawatt  Unit = "megawatt"
	UnitGigawatt  Unit = "gigawatt"
	UnitTerrawatt Unit = "terrawatt"

	// CURRENT
	UnitMilliampere Unit = "milliampere"
	UnitAmpere      Unit = "ampere"

	// POTENTIAL
	UnitMillivolt Unit = "millivolt"
	UnitVolt      Unit = "volt"

	// APM
	UnitSpan Unit = "span"

	// SYNTHETICS
	UnitRun Unit = "run"
)
