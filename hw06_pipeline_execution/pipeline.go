package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	I   = interface{}
	In  = <-chan I
	Out = In
	Bi  = chan I
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	mydone := make(chan struct{}) //придется делать свой done, чтобы его закрыть широковещательно
	go func() {
		<-done        //к сожалению в тестах закрывается только для 1 читателя
		close(mydone) //закрыть широковещательно
	}()
	out := in
	for _, stage := range stages {
		newIn := make(Bi) //делаем свой входящий BI-канал каждый раз,т.к. я не могу закрывать каналы только для чтнения
		go func(_in Bi, _out Out) {
			defer close(_in) //схлопываем каналы
			for {
				select {
				case <-mydone:
					return
				case v, ok := <-_out: //default не подходит-т.к. не блокирует мулльтиплексор
					if !ok {
						return
					}
					_in <- v //перекачиваем данные с выхода (только для чтения) на вход Bi
				}
			}
		}(newIn, out)
		out = stage(newIn)
	}
	return out
}
